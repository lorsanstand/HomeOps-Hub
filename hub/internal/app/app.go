package app

import (
	"database/sql"
	"fmt"
	standartlog "log"
	"net"
	"os"

	hubdir "github.com/lorsanstand/HomeOps-Hub/hub/internal"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/migrator"
	grpcserv "github.com/lorsanstand/HomeOps-Hub/hub/internal/rpc"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/service/connection_manager"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/service/hub_service"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/store"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/utils/notifier"
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
		os.Exit(1)
	}

	defer func() {
		if err := DBConn.Close(); err != nil {
			a.log.Warn().Err(err).Msg("failed to close migrate postgres connection")
		}
	}()

	a.log.Info().Msg("applying database migrations")
	if err := applyMigrations(DBConn); err != nil {
		a.log.Error().Err(err).Msg("")
		os.Exit(1)
	}
	a.log.Info().Msg("migrations applied successfully")

	hubStore := store.NewHubStore(DBConn)
	hubService := hub_service.NewHubService(hubStore, a.log)
	statusNotifier := notifier.NewStatusNotifier()
	connManger := connection_manager.NewConnectionManager(hubStore, statusNotifier, a.log)

	a.log.Info().Msg("starting hub service")

	err = a.hubServe(hubService, connManger)
	if err != nil {
		a.log.Error().Err(err).Msg("hub service failed to start")
		return
	}
}

func (a *App) hubServe(hubService *hub_service.HubService, manager *connection_manager.ConnectionManager) error {
	address := fmt.Sprintf("0.0.0.0:%v", a.cfg.Port)

	server := grpcserv.NewHubHandler(hubService, manager, a.log)

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

func applyMigrations(db *sql.DB) error {
	mgrt, err := migrator.NewMigrator(hubdir.MigrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	if err = mgrt.ApplyMigrations(db); err != nil {
		return fmt.Errorf("migrations failed to apply: %w", err)
	}
	return nil
}
