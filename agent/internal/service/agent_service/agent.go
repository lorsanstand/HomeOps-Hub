package agent_service

import (
	"context"
	"fmt"

	"github.com/lorsanstand/HomeOps-Hub/agent/internal/utils/config_yaml"
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/utils/settings"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
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
	logger = logger.With().Str("component", "internal.service.agent_serivce").Logger()

	return &AgentService{collect: collector, conn: conn, cfg: cfg, log: logger, settings: settings}
}

func (a *AgentService) RegisterAgentConn(ctx context.Context) error {
	a.log.Debug().Msg("getting info by system")
	info, caps := a.collect.GatherInfoSystem()
	a.log.Debug().Msg("create request data for register agent")
	AgentID := a.settings.AgentID
	AgentName := a.cfg.AppName
	AgentData := domain.RegisterAgentRequest{
		AgentId:      AgentID,
		AgentName:    AgentName,
		Host:         info,
		Capabilities: caps,
		AgentVersion: a.cfg.GetAgentVersion(),
	}

	data, err := a.conn.RegisterAgent(ctx, AgentData)
	if err != nil {
		return fmt.Errorf("register agent: %w", err)
	}

	if err = a.settings.Insert(settings.Settings{AgentID: data.AgentID}); err != nil {
		return fmt.Errorf("save agent ID: %w", err)
	}

	return nil
}
