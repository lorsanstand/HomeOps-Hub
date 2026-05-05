package hub_service

import (
	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

func toCreateAgentModel(agent domain.RegisterAgentRequest) domainHub.CreateAgentModel {
	return domainHub.CreateAgentModel{
		AgentID:      agent.AgentID,
		AgentName:    agent.AgentName,
		Architecture: agent.Host.Arch,
		System:       agent.Host.System,
		Hostname:     agent.Host.Hostname,
		Version:      agent.AgentVersion,
		Capabilities: agent.Capabilities,
	}
}
