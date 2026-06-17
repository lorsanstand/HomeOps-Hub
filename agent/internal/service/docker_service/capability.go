package docker_service

import domainHub "github.com/lorsanstand/HomeOps-Hub/shared/domain"

func (d *DockerService) Capability() domainHub.Capability {
	capability := domainHub.Capability{
		Available: true,
		Version:   "0",
		Name:      "docker",
		Reason:    "",
	}

	for command := com
}
