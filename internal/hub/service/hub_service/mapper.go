package hub_service

import (
	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	domainHub "github.com/lorsanstand/HomeOps-Hub/internal/hub/domain"
)

func toCreateAgentModel(agent domain.RegisterAgentRequest) domainHub.CreateAgentModel {
	return domainHub.CreateAgentModel{
		AgentID:      agent.AgentId,
		AgentName:    agent.AgentName,
		Architecture: agent.Host.Arch,
		System:       agent.Host.System,
		Hostname:     agent.Host.Hostname,
		Version:      agent.AgentVersion,
		Capabilities: agent.Capabilities,
	}
}
