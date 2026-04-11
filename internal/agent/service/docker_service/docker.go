package docker_service

import (
	"context"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
	"github.com/rs/zerolog"
)

type dockerAPI interface {
	Ping(ctx context.Context) (types.Ping, error)
	ContainerList(ctx context.Context, opts container.ListOptions) ([]container.Summary, error)
}

type DockerService struct {
	dockerClient dockerAPI
	log          zerolog.Logger
}

func NewDockerService(api dockerAPI, logger zerolog.Logger) *DockerService {
	return &DockerService{
		dockerClient: api,
		log:          logger.With().Str("component", "agent.serivce.docker").Logger(),
	}
}

func (d *DockerService) CheckDockerDaemon(ctx context.Context) error {
	_, err := d.dockerClient.Ping(ctx)
	d.log.Debug().Msg("check docker")
	return err
}

func (d *DockerService) ContainersList(ctx context.Context) ([]container.Summary, error) {
	ContainersList, err := d.dockerClient.ContainerList(ctx, container.ListOptions{})
	d.log.Debug().Msg("get container list")
	return ContainersList, err
}
