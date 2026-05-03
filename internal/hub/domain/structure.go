package domain

import (
	"time"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
)

type CreateAgentModel struct {
	AgentID      string
	AgentName    string
	Architecture string
	System       string
	Hostname     string
	Version      string
	Capabilities []domain.Capability
}

type AgentModel struct {
	ID           int
	AgentID      string
	AgentName    string
	Architecture string
	System       string
	Hostname     string
	Version      string
	Capabilities []domain.Capability
	RegisteredAt time.Time
}
