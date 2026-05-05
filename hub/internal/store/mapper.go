package store

import (
	"encoding/json"

	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	gen2 "github.com/lorsanstand/HomeOps-Hub/hub/internal/store/sqlc/gen"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

func toDBAgent(agent domainHub.CreateAgentModel) gen2.CreateAgentParams {
	jsonCaps := toJSONCapabilities(agent.Capabilities)
	capsStr := string(jsonCaps)
	return gen2.CreateAgentParams{
		AgentID:      agent.AgentID,
		AgentName:    &agent.AgentName,
		Architecture: agent.Architecture,
		System:       agent.System,
		Hostname:     agent.Hostname,
		Version:      agent.Version,
		Capabilities: &capsStr,
	}
}

func toUpdateDBAgent(agent domainHub.CreateAgentModel) gen2.UpdateAgentByIDParams {
	jsonCaps := toJSONCapabilities(agent.Capabilities)
	capsStr := string(jsonCaps)
	return gen2.UpdateAgentByIDParams{
		AgentID:      agent.AgentID,
		AgentName:    &agent.AgentName,
		Architecture: agent.Architecture,
		System:       agent.System,
		Hostname:     agent.Hostname,
		Version:      agent.Version,
		Capabilities: &capsStr,
	}
}

func toJSONCapabilities(caps []domain.Capability) []byte {
	data, err := json.Marshal(caps)
	if err != nil {
		// Note: Error is silently handled - consider logging in production
		return []byte{}
	}
	return data
}

func toAgentModel(dbAgent gen2.Agent) domainHub.AgentModel {
	var dbAgentName string
	if dbAgent.AgentName != nil {
		dbAgentName = *dbAgent.AgentName
	}

	var capsBytes []byte
	if dbAgent.Capabilities != nil {
		capsBytes = []byte(*dbAgent.Capabilities)
	}

	return domainHub.AgentModel{
		ID:           int(dbAgent.ID),
		AgentID:      dbAgent.AgentID,
		AgentName:    dbAgentName,
		Architecture: dbAgent.Architecture,
		System:       dbAgent.System,
		Hostname:     dbAgent.Hostname,
		Capabilities: toDomainCapabilities(capsBytes),
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
