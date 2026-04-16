package store

import (
	"encoding/json"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/store/sqlc/gen"
)

type Agent struct {
	AgentID      string
	AgentName    string
	Architecture string
	System       string
	Hostname     string
	Version      string
	Capabilities []domain.Capability
}

func toDBAgent(agent Agent) gen.CreateAgentParams {
	return gen.CreateAgentParams{
		AgentID:      agent.AgentID,
		AgentName:    &agent.AgentName,
		Architecture: agent.Architecture,
		System:       agent.System,
		Version:      agent.Version,
		Capabilities: toJsonCapabilities(agent.Capabilities),
	}
}

func toJsonCapabilities(caps []domain.Capability) []byte {
	data, err := json.Marshal(caps)
	if err != nil {
		return []byte{}
	}
	return data
}
