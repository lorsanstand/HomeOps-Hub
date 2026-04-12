package app

import (
	standartlog "log"

	"github.com/docker/docker/client"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/rpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/service/collector"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/service/docker_service"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/config_yaml"
	log2 "github.com/lorsanstand/HomeOps-Hub/internal/shared/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type App struct {
	log     zerolog.Logger
	cfg     *config_yaml.AgentConfig
	hubConn *rpc.Connection
}

func NewApp() *App {
	cfg, err := config_yaml.NewConfig()
	if err != nil {
		standartlog.Fatalf("failed get config: %v", err)
	}

	log := log2.NewLogger(cfg)

	return &App{cfg: cfg, log: log}
}

func (a *App) Run() {
	GRPCConn, err := grpc.NewClient(a.cfg.GetGRPCAddress())
	if err != nil {
		a.log.Error().Err(err).Msg("failed to get hub connections")
		return
	}

	conn := rpc.NewConnectAgent(GRPCConn, a.log)
	defer conn.Close()

	var DockerService collector.Docker

	DockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		a.log.Warn().Err(err).Msg("failed to get docker API")
		DockerService = docker_service.NewBadDocker("not_installed")
	} else {
		DockerService = docker_service.NewDockerService(DockerClient, a.log)
	}

	Collector := collector.NewCollector(DockerService, a.log)

	Collector.GatherInfoSystem()
}
