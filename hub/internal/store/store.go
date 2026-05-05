package store

import (
	"context"
	"database/sql"

	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/store/sqlc/gen"
)

type HubStore struct {
	queries *gen.Queries
}

func NewHubStore(db *sql.DB) *HubStore {
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
	data.ID = int64(ID)
	return h.queries.UpdateAgentByID(ctx, data)
}
