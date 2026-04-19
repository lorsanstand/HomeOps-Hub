package app

import (
	"context"
	"database/sql"
	"fmt"
	standartlog "log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	hubdir "github.com/lorsanstand/HomeOps-Hub/internal/hub"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/migrator"
	grpcserv "github.com/lorsanstand/HomeOps-Hub/internal/hub/rpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/service/hub_service"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/store"
	"github.com/lorsanstand/HomeOps-Hub/internal/shared/config"
	"github.com/lorsanstand/HomeOps-Hub/internal/shared/log"
	"github.com/rs/zerolog"
)

type App struct {
	cfg *config.Config
	log zerolog.Logger
}

func NewApp() *App {
	cfg, err := config.NewConfig()
	if err != nil {
		standartlog.Fatalf("failed get config: %v", err)
	}

	logger := log.NewLogger(cfg)

	return &App{cfg: cfg, log: logger}
}

func (a *App) Run() {
	ctx := context.Background()
	a.log.Info().Str("host", a.cfg.DBHost).Int("port", a.cfg.DBPort).Msg("connecting to database")
	migratePGConn, err := sql.Open("pgx", a.cfg.GetURLPostgres())
	if err != nil {
		a.log.Error().Err(err).Msg("failed to connect to the database for migrations")
		return
	}
	defer migratePGConn.Close()

	mgrt, err := migrator.NewMigrator(hubdir.MigrationsFS, "migrations")
	if err != nil {
		a.log.Error().Err(err).Msg("failed to create migrator")
		return
	}

	a.log.Info().Msg("applying database migrations")
	if err = mgrt.ApplyMigrations(migratePGConn); err != nil {
		a.log.Error().Err(err).Msg("migrations failed to apply")
		return
	}
	a.log.Info().Msg("migrations applied successfully")
	migratePGConn.Close()

	a.log.Info().Msg("creating database connection pool")
	pool, err := pgxpool.New(ctx, a.cfg.GetURLPostgres())
	if err != nil {
		a.log.Error().Err(err).Msg("failed to create database connection pool")
		return
	}
	defer pool.Close()
	a.log.Info().Msg("database connection pool created")

	hubStore := store.NewHubStore(pool)
	hubService := hub_service.NewHubService(hubStore, a.log)

	a.log.Info().Msg("starting hub service")
	err = a.hubServe(hubService)
	if err != nil {
		a.log.Error().Err(err).Msg("hub service failed to start")
		return
	}
}

func (a *App) hubServe(hubService *hub_service.HubService) error {
	address := fmt.Sprintf("0.0.0.0:%v", a.cfg.Port)
	a.log.Info().Str("address", address).Msg("starting gRPC server")

	server := grpcserv.NewHubHandler(hubService, a.log)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		a.log.Error().Err(err).Str("address", address).Msg("failed to listen on address")
		return err
	}
	a.log.Info().Str("address", address).Msg("listening on address")

	a.log.Info().Msg("gRPC server is running")
	err = server.GrpcServer.Serve(lis)
	if err != nil {
		a.log.Error().Err(err).Msg("gRPC server error")
		return err
	}

	return nil
}
