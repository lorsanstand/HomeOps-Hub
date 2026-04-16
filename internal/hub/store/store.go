package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/store/sqlc/gen"
)

type HubStore struct {
	queries *gen.Queries
}

func NewHubStore(db *pgxpool.Pool) *HubStore {
	queries := gen.New(db)
	return &HubStore{queries}
}

func (h *HubStore) NewAgent(ctx context.Context, agent Agent) error {
	return h.queries.CreateAgent(ctx, toDBAgent(agent))
}
