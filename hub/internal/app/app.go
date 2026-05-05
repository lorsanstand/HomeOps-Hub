package app

import (
	"database/sql"
	"fmt"
	standartlog "log"
	"net"

	hubdir "github.com/lorsanstand/HomeOps-Hub/hub/internal"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/migrator"
	grpcserv "github.com/lorsanstand/HomeOps-Hub/hub/internal/rpc"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/service/hub_service"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/store"
	"github.com/lorsanstand/HomeOps-Hub/shared/config"
	"github.com/lorsanstand/HomeOps-Hub/shared/log"
	_ "github.com/mattn/go-sqlite3"
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
		return nil
	}

	logger := log.NewLogger(cfg)

	return &App{cfg: cfg, log: logger}
}

func (a *App) Run() {
	a.log.Info().Msg("connecting to database")
	DBConn, err := sql.Open("sqlite", "database.db")
	if err != nil {
		a.log.Error().Err(err).Msg("failed to connect to the database")
		return
	}

	defer func() {
		if err := DBConn.Close(); err != nil {
			a.log.Warn().Err(err).Msg("failed to close migrate postgres connection")
		}
	}()

	mgrt, err := migrator.NewMigrator(hubdir.MigrationsFS, "migrations")
	if err != nil {
		a.log.Error().Err(err).Msg("failed to create migrator")
		return
	}

	a.log.Info().Msg("applying database migrations")
	if err = mgrt.ApplyMigrations(DBConn); err != nil {
		a.log.Error().Err(err).Msg("migrations failed to apply")
		return
	}
	a.log.Info().Msg("migrations applied successfully")

	hubStore := store.NewHubStore(DBConn)
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
