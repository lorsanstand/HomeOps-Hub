package app

import (
	"fmt"
	standartlog "log"
	"net"

	grpcserv "github.com/lorsanstand/HomeOps-Hub/internal/hub/rpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/service/hub_service"
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
	err := a.hubServe()
	if err != nil {
		a.log.Error().Err(err).Msg("failed to start the server")
	}
}

func (a *App) hubServe() error {
	address := fmt.Sprintf("0.0.0.0:%v", a.cfg.Port)
	a.log.Info().Str("address", "http://"+address).Msg("start GRPC server")

	hub := hub_service.NewHubService(a.log)

	server := grpcserv.NewHubHandler(hub, a.log)

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
