package collector

import (
	"os"
	"runtime"

	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
	"github.com/rs/zerolog"
)

type Docker interface {
	Capability() domain.Capability
}

type Collector struct {
	log          zerolog.Logger
	dockerReader Docker
}

func NewCollector(docker Docker, logger zerolog.Logger) *Collector {
	logger = logger.With().Str("component", "cmd.service.collector").Logger()

	return &Collector{log: logger, dockerReader: docker}
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

	caps := []domain.Capability{c.dockerReader.Capability()}
	return host, caps
}
