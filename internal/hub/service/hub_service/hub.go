package hub_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	domainHub "github.com/lorsanstand/HomeOps-Hub/internal/hub/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/utils/hasher"
	"github.com/rs/zerolog"
)

type Store interface {
	NewAgent(ctx context.Context, agent domainHub.CreateAgentModel) error
	GetAgentByAgentID(ctx context.Context, AgentID string) (domainHub.AgentModel, error)
	UpdateAgentByID(ctx context.Context, ID int, updateAgent domainHub.CreateAgentModel) error
}

type HubService struct {
	store Store
	log   zerolog.Logger
}

func NewHubService(store Store, logger zerolog.Logger) *HubService {
	return &HubService{log: logger, store: store}
}

func (h *HubService) RegisterAgent(ctx context.Context, data domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error) {
	h.log.Debug().Str("agentId", data.AgentId).Str("agentName", data.AgentName).Msg("started registering agent")
	agent, err := h.store.GetAgentByAgentID(ctx, data.AgentId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.log.Error().Err(err).Str("agentId", data.AgentId).Msg("failed to get agent from database")
		return domain.RegisterAgentResponse{}, fmt.Errorf("failed select agent to db: %w", err)
	}

	if data.AgentId != "" && !errors.Is(err, sql.ErrNoRows) {
		h.log.Debug().Str("agentId", agent.AgentID).Str("agentName", data.AgentName).Msg("agent exists, updating")

		data.AgentId = agent.AgentID

		agentStore := toCreateAgentModel(data)

		if err := h.store.UpdateAgentByID(ctx, agent.ID, agentStore); err != nil {
			h.log.Error().Err(err).Str("agentId", agent.AgentID).Msg("failed to update agent in database")
			return domain.RegisterAgentResponse{}, err
		}
		h.log.Info().Str("agentId", agent.AgentID).Msg("agent updated successfully")
		return domain.RegisterAgentResponse{AgentID: agent.AgentID, Heartbeat: 5}, nil
	}

	AgentID, err := hasher.MakeID(data.Host, data.AgentName)
	if err != nil {
		h.log.Error().Err(err).Str("agentName", data.AgentName).Str("hostname", data.Host.Hostname).Msg("failed to generate agent id")
		return domain.RegisterAgentResponse{}, err
	}

	data.AgentId = AgentID

	agentStore := toCreateAgentModel(data)

	if err := h.store.NewAgent(ctx, agentStore); err != nil {
		h.log.Error().Err(err).Str("agentId", AgentID).Str("agentName", data.AgentName).Msg("failed to create new agent in database")
		return domain.RegisterAgentResponse{}, err
	}

	h.log.Info().Str("agentId", AgentID).Str("agentName", data.AgentName).Str("hostname", data.Host.Hostname).Msg("agent registered successfully")
	return domain.RegisterAgentResponse{AgentID: AgentID, Heartbeat: 5}, nil
}
