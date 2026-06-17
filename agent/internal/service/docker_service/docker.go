package docker_service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
		log:          logger,
	}
}

func (d *DockerService) CheckDockerDaemon(ctx context.Context, args map[string]string) (string, error) {
	_, err := d.dockerClient.Ping(ctx)
	d.log.Debug().Msg("ping docker")
	if err != nil {
		return "", err
	}
	return "successful", nil
}

func (d *DockerService) ContainersList(ctx context.Context) ([]container.Summary, error) {
	ContainersList, err := d.dockerClient.ContainerList(ctx, container.ListOptions{})
	d.log.Debug().Msg("get container list")
	return ContainersList, err
}
