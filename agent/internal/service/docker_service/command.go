package docker_service

import (
	"github.com/lorsanstand/HomeOps-Hub/agent/internal/domain"
	domainHub "github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

func (d *DockerService) commands() []domain.Command {
	return []domain.Command{
		{Execute: d.CheckDockerDaemon, CapabilityCommand: domainHub.CapabilityCommand{
			Name:         "check",
			OptionalArgs: []domainHub.CommandArgs{},
			RequiredArgs: []domainHub.CommandArgs{},
			Version:      "0",
			Description:  "Check docker daemon",
			TypeOutput: "string",
		}},
	}
}
