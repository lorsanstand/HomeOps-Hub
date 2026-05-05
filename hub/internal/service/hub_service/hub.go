package hub_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/utils/hasher"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
	"github.com/rs/zerolog"
)

const HEARTBEAT = 5

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
	h.log.Debug().Str("agentID", data.AgentID).Str("agentName", data.AgentName).Msg("started registering agent")
	agent, err := h.store.GetAgentByAgentID(ctx, data.AgentID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return domain.RegisterAgentResponse{}, fmt.Errorf("failed select agent to db: %w", err)
	}

	if data.AgentID != "" && !errors.Is(err, sql.ErrNoRows) {
		h.log.Debug().Str("agentID", agent.AgentID).Str("agentName", data.AgentName).Msg("agent exists, updating")

		data.AgentID = agent.AgentID

		agentStore := toCreateAgentModel(data)

		if err := h.store.UpdateAgentByID(ctx, agent.ID, agentStore); err != nil {
			return domain.RegisterAgentResponse{}, fmt.Errorf("update agent in db: %w", err)
		}
		h.log.Debug().Str("agentId", agent.AgentID).Msg("agent updated successfully")
		return domain.RegisterAgentResponse{AgentID: agent.AgentID, Heartbeat: HEARTBEAT}, nil
	}

	AgentID, err := hasher.MakeID(data.Host, data.AgentName)
	if err != nil {
		return domain.RegisterAgentResponse{}, fmt.Errorf("generate agent ID: %w", err)
	}

	data.AgentID = AgentID

	agentStore := toCreateAgentModel(data)

	if err := h.store.NewAgent(ctx, agentStore); err != nil {
		return domain.RegisterAgentResponse{}, fmt.Errorf("insert new agent: %w", err)
	}
	return domain.RegisterAgentResponse{AgentID: AgentID, Heartbeat: HEARTBEAT}, nil
}
