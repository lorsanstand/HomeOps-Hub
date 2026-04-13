package agent_service

import (
	"context"
	"fmt"

	"github.com/lorsanstand/HomeOps-Hub/internal/agent/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/utils/config_yaml"
	"github.com/rs/zerolog"
)

type Collector interface {
	GatherInfoSystem() (domain.HostInfo, []domain.Capability)
}

type HubConnection interface {
	RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentData) (domain.RegisterAgentDataResponse, error)
}

type AgentService struct {
	collect   Collector
	conn      HubConnection
	log       zerolog.Logger
	cfg       *config_yaml.AgentConfig
	heartBeat int
	agentID   string
}

func NewAgentService(
	collector Collector,
	conn HubConnection,
	AgentID string,
	cfg *config_yaml.AgentConfig,
	logger zerolog.Logger,
) *AgentService {
	logger = logger.With().Str("component", "agent.service.agent_serivce").Logger()

	return &AgentService{collect: collector, conn: conn, cfg: cfg, log: logger, agentID: AgentID}
}

func (a *AgentService) RegisterAgentConn(ctx context.Context) {
	info, caps := a.collect.GatherInfoSystem()
	AgentID := a.agentID
	AgentName := a.cfg.AppName
	AgentData := domain.RegisterAgentData{AgentId: AgentID, AgentName: AgentName, Host: info, Capabilities: caps}

	data, err := a.conn.RegisterAgent(ctx, AgentData)
	if err != nil {
		a.log.Error().Err(err).Msg("failed register agent")
		return
	}
	fmt.Println(data)
}
