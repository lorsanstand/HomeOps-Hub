package app

import (
	"context"
	standartlog "log"

	"github.com/docker/docker/client"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/rpc"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/service/agent_service"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/service/collector"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/service/docker_service"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/utils/config_yaml"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/utils/settings"
	"github.com/lorsanstand/HomeOps-Hub/shared/log"
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

	logger := log.NewLogger(cfg)
	logger = logger.With().Str("component", "internal.app").Logger()
	logger = logger.With().Str("name", cfg.AppName).Logger()

	sett, err := settings.ReadSettings(cfg.SettingsPath)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get settings")
		return nil, err
	}

	return &App{cfg: cfg, log: logger, settings: sett}, nil
}

func (a *App) Run() {
	ctx := context.Background()

	GRPCConn, err := grpc.NewClient(a.cfg.GetGRPCAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		a.log.Error().Err(err).Msg("failed to connection hub")
		return
	}
	a.log.Info().Msg("connection to the hub successful")

	conn := rpc.NewConnectAgent(GRPCConn)
	defer func() {
		if err := conn.Close(); err != nil {
			a.log.Warn().Err(err).Msg("failed to close rpc connection")
		}
	}()

	var DockerService collector.Docker

	DockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		a.log.Warn().Err(err).Msg("failed to get docker API")
		DockerService = docker_service.NewBadDocker("not_installed")
	} else {
		a.log.Info().Msg("successfully to get docker API")
		DockerService = docker_service.NewDockerService(DockerClient, a.log)
	}

	collect := collector.NewCollector(DockerService, a.log)

	agent := agent_service.NewAgentService(collect, conn, a.settings, a.cfg, a.log)
	if err := agent.RegisterAgentConn(ctx); err != nil {
		a.log.Error().Err(err).Msg("failed to agent registration")
	}
	a.log.Info().Msg("agent registration complete")
}
