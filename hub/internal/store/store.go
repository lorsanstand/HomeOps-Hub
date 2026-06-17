package store

import (
	"context"
	"database/sql"
	"time"

	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/store/sqlc/gen"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
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

func (h *HubStore) GetAgentByAgentID(ctx context.Context, agentID string) (domainHub.AgentModel, error) {
	data, err := h.queries.GetAgentByAgentID(ctx, agentID)
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

func (h *HubStore) CreateHeartbeat(ctx context.Context, heartbeat domain.CreateHeartbeatModel) error {
	data := toDBHeartbeat(heartbeat)
	return h.queries.InsertHeartbeat(ctx, data)
}

func (h *HubStore) GetHeartbeatsByIDAfter(ctx context.Context, agentID string, timestamp time.Time) ([]domain.HeartbeatModel, error) {
	data := gen.SelectHeartbeatsAfterParams{AgentID: agentID, Timestamp: timestamp}
	heartbeats, err := h.queries.SelectHeartbeatsAfter(ctx, data)
	if err != nil {
		return []domain.HeartbeatModel{}, err
	}

	heartbeatsModel := make([]domain.HeartbeatModel, len(heartbeats))

	for i, heartbeat := range heartbeats {
		heartbeatsModel[i] = toHeartBeatModel(heartbeat)
	}

	return heartbeatsModel, nil
}
