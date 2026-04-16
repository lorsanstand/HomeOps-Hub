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
	migratePGConn, err := sql.Open("pgx", a.cfg.GetURLPostgres())
	if err != nil {
		a.log.Error().Err(err).Msg("failed to connect to the database")
		return
	}
	defer migratePGConn.Close()

	mgrt, err := migrator.NewMigrator(hubdir.MigrationsFS, "migrations")
	if err != nil {
		a.log.Error().Err(err).Msg("failed create migrator")
		return
	}

	if err = mgrt.ApplyMigrations(migratePGConn); err != nil {
		a.log.Error().Err(err).Msg("migrations were not applied")
		return
	}
	migratePGConn.Close()

	pool, err := pgxpool.New(ctx, a.cfg.GetURLPostgres())
	if err != nil {
		a.log.Error().Err(err).Msg("failed create db pool")
		return
	}

	hubStore := store.NewHubStore(pool)

	hubService := hub_service.NewHubService(hubStore, a.log)

	err = a.hubServe(hubService)
	if err != nil {
		a.log.Error().Err(err).Msg("failed to start the server")
		return
	}
}

func (a *App) hubServe(hubService *hub_service.HubService) error {
	address := fmt.Sprintf("0.0.0.0:%v", a.cfg.Port)
	a.log.Info().Str("address", "http://"+address).Msg("start GRPC server")

	server := grpcserv.NewHubHandler(hubService, a.log)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = server.GrpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
