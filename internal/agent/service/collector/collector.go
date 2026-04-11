package collector

import (
	"context"
	"os"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/domain"
	"github.com/rs/zerolog"
)

type Docker interface {
	Ping(ctx context.Context) (types.Ping, error)
	ContainerList(ctx context.Context, opts container.ListOptions) ([]container.Summary, error)
	Capability() domain.Capability
}

type Collector struct {
	log    zerolog.Logger
	docker Docker
}

func NewCollector(docker Docker, logger zerolog.Logger) *Collector {
	logger = logger.With().Str("component", "agent.service.collector").Logger()

	return &Collector{log: logger, docker: docker}
}

func (c *Collector) GatherInfoSystem() (domain.HostInfo, []domain.Capability) {
	var host domain.HostInfo

	hostname, err := os.Hostname()
	if err != nil {
		c.log.Warn().Msg("failed to get hostname")
	}
	host.Hostname = hostname

	host.Arch = runtime.GOARCH
	host.System = runtime.GOOS

	caps := []domain.Capability{c.docker.Capability()}
	return host, caps
}
