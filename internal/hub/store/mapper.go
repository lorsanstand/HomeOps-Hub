package store

import (
	"encoding/json"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	domainHub "github.com/lorsanstand/HomeOps-Hub/internal/hub/domain"
	"github.com/lorsanstand/HomeOps-Hub/internal/hub/store/sqlc/gen"
)

func toDBAgent(agent domainHub.CreateAgentModel) gen.CreateAgentParams {
	return gen.CreateAgentParams{
		AgentID:      agent.AgentID,
		AgentName:    &agent.AgentName,
		Architecture: agent.Architecture,
		System:       agent.System,
		Hostname:     agent.Hostname,
		Version:      agent.Version,
		Capabilities: toJsonCapabilities(agent.Capabilities),
	}
}

func toUpdateDBAgent(agent domainHub.CreateAgentModel) gen.UpdateAgentByIDParams {
	return gen.UpdateAgentByIDParams{
		AgentID:      agent.AgentID,
		AgentName:    &agent.AgentName,
		Architecture: agent.Architecture,
		System:       agent.System,
		Hostname:     agent.Hostname,
		Version:      agent.Version,
		Capabilities: toJsonCapabilities(agent.Capabilities),
	}
}

func toJsonCapabilities(caps []domain.Capability) []byte {
	data, err := json.Marshal(caps)
	if err != nil {
		// Note: Error is silently handled - consider logging in production
		return []byte{}
	}
	return data
}

func toAgentModel(dbAgent gen.Agent) domainHub.AgentModel {
	var dbAgentName string
	if dbAgent.AgentName != nil {
		dbAgentName = *dbAgent.AgentName
	}

	return domainHub.AgentModel{
		ID:           int(dbAgent.ID),
		AgentID:      dbAgent.AgentID,
		AgentName:    dbAgentName,
		Architecture: dbAgent.Architecture,
		System:       dbAgent.System,
		Hostname:     dbAgent.Hostname,
		Capabilities: toDomainCapabilities(dbAgent.Capabilities),
	}
}

func toDomainCapabilities(caps []byte) []domain.Capability {
	var capabilities []domain.Capability
	err := json.Unmarshal(caps, &capabilities)
	if err != nil {
		// Note: Error is silently handled - consider logging in production
		return []domain.Capability{}
	}
	return capabilities
}
