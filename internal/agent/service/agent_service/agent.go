package agent_service

import (
	"context"
	"fmt"

	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/config_yaml"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/settings"
	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/rs/zerolog"
)

type Collector interface {
	GatherInfoSystem() (domain.HostInfo, []domain.Capability)
}

type HubConnection interface {
	RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error)
}

type AgentService struct {
	collect   Collector
	conn      HubConnection
	log       zerolog.Logger
	cfg       *config_yaml.AgentConfig
	heartBeat int
	settings  *settings.Settings
}

func NewAgentService(
	collector Collector,
	conn HubConnection,
	settings *settings.Settings,
	cfg *config_yaml.AgentConfig,
	logger zerolog.Logger,
) *AgentService {
	logger = logger.With().Str("component", "agent.service.agent_serivce").Logger()

	return &AgentService{collect: collector, conn: conn, cfg: cfg, log: logger, settings: settings}
}

func (a *AgentService) RegisterAgentConn(ctx context.Context) {
	info, caps := a.collect.GatherInfoSystem()
	AgentID := a.settings.AgentID
	AgentName := a.cfg.AppName
	AgentData := domain.RegisterAgentRequest{AgentId: AgentID, AgentName: AgentName, Host: info, Capabilities: caps}

	data, err := a.conn.RegisterAgent(ctx, AgentData)
	if err != nil {
		a.log.Error().Err(err).Msg("failed register agent")
		return
	}

	if err = a.settings.Insert(settings.Settings{AgentID: data.AgentID}); err != nil {
		a.log.Warn().Err(err).Msg("failed to save agent id")
	}
	fmt.Println(data)
}
