package hub_service

import (
	"context"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/store"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/utils/hasher"
	"github.com/rs/zerolog"
)

type Store interface {
	NewAgent(ctx context.Context, agent store.Agent) error
}

type HubService struct {
	store Store
	log   zerolog.Logger
}

func NewHubService(store Store, logger zerolog.Logger) *HubService {
	return &HubService{log: logger, store: store}
}

func (h *HubService) RegisterAgent(ctx context.Context, data domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error) {
	h.log.Debug().Msg("registered new agent")
	AgentID := data.AgentId
	if data.AgentId == "" {
		var err error
		AgentID, err = hasher.MakeID(data.Host, data.AgentName)
		if err != nil {
			h.log.Error().Err(err).Msg("failed create agent id")
			AgentID = ""
		}
	}

	agentStore := store.Agent{
		AgentID:      AgentID,
		AgentName:    data.AgentName,
		Architecture: data.Host.Arch,
		System:       data.Host.System,
		Hostname:     data.Host.Hostname,
		Version:      data.AgentVersion,
		Capabilities: data.Capabilities,
	}

	if err := h.store.NewAgent(ctx, agentStore); err != nil {
		h.log.Warn().Err(err).Msg("failed add new agent in db")
		return domain.RegisterAgentResponse{}, err
	}

	return domain.RegisterAgentResponse{AgentID: AgentID, Heartbeat: 5}, nil
}
