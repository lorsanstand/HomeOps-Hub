package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	domainHub "github.com/lorsanstand/HomeOps-Hub/internal/hub/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/store/sqlc/gen"
)

type HubStore struct {
	queries *gen.Queries
}

func NewHubStore(db *pgxpool.Pool) *HubStore {
	queries := gen.New(db)
	return &HubStore{queries}
}

func (h *HubStore) NewAgent(ctx context.Context, agent domainHub.CreateAgentModel) error {
	return h.queries.CreateAgent(ctx, toDBAgent(agent))
}

func (h *HubStore) GetAgentByAgentID(ctx context.Context, AgentID string) (domainHub.AgentModel, error) {
	data, err := h.queries.GetAgentByAgentID(ctx, AgentID)
	if err != nil {
		return domainHub.AgentModel{}, err
	}
	return toAgentModel(data), nil
}

func (h *HubStore) UpdateAgentByID(ctx context.Context, ID int, updateAgent domainHub.CreateAgentModel) error {
	data := toUpdateDBAgent(updateAgent)
	data.ID = int32(ID)
	return h.queries.UpdateAgentByID(ctx, data)
}
