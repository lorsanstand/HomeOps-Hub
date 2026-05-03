package app

import (
	"context"
	standartlog "log"

	"github.com/docker/docker/client"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/rpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/service/agent_service"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/service/collector"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/service/docker_service"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/config_yaml"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/settings"
	log2 "github.com/lorsanstand/HomeOps-Hub/internal/shared/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log      zerolog.Logger
	cfg      *config_yaml.AgentConfig
	settings *settings.Settings
	hubConn  *rpc.Connection
}

func NewApp() (*App, error) {

	cfg, err := config_yaml.NewConfig()
	if err != nil {
		standartlog.Fatalf("failed to get config: %v", err)
		return nil, err
	}

	log := log2.NewLogger(cfg)
	log = log.With().Str("component", "agent.app").Logger()

	sett, err := settings.ReadSettings(cfg.SettingsPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to get settings")
		return nil, err
	}

	return &App{cfg: cfg, log: log, settings: sett}, nil
}

func (a *App) Run() {
	ctx := context.Background()

	GRPCConn, err := grpc.NewClient(a.cfg.GetGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	collect := collector.NewCollector(DockerService, a.log)

	agent := agent_service.NewAgentService(collect, conn, a.settings, a.cfg, a.log)
	agent.RegisterAgentConn(ctx)
}
