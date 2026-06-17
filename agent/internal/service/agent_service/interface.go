package agent_service

import (
	"context"

	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

type Collector interface {
	GatherInfoSystem() (domain.HostInfo, []domain.Capability)
}

type Settings interface {
	InsertAgentID(agentID string) error
	GetAgentID() string
}

type HubConnection interface {
	RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error)
}
