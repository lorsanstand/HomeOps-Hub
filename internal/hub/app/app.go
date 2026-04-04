package app

import (
	"fmt"
	"net"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	grpcserv "github.com/lorsanstand/HomeOps-Hub/internal/hub/grpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/shared/config"
	"github.com/lorsanstand/HomeOps-Hub/internal/shared/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type App struct {
	cfg *config.Config
	log zerolog.Logger
}

func NewApp() *App {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Errorf("failed get config: %v", err)
	}

	logger := log.NewLogger(cfg)

	return &App{cfg: cfg, log: logger}
}

func (a *App) Run() {
	address := fmt.Sprintf("http://0.0.0.0:%v", a.cfg.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.cfg.Port))
	if err != nil {
		a.log.Error().Err(err).Msg("failed started listen")
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHubServer(grpcServer, grpcserv.NewServer())

	a.log.Info().Str("address", address).Msg("server started")

	err = grpcServer.Serve(lis)
	if err != nil {
		a.log.Error().Err(err).Msg("failed started grpc server")
		return
	}

}
