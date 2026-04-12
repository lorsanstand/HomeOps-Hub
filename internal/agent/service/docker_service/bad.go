package docker_service

import "github.com/lorsanstand/HomeOps-Hub/internal/agent/domain"

type BadDocker struct {
	reason string
}

func (d *BadDocker) Capability() domain.Capability {
	return domain.Capability{
		Name:      "docker",
		Available: false,
		Version:   "",
		Reason:    d.reason,
	}
}

func NewBadDocker(reason string) *BadDocker {
	return &BadDocker{reason: reason}
}
