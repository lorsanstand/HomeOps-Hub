package docker_service

import (
	"context"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
)

type dockerAPI interface {
	Ping(ctx context.Context) (types.Ping, error)
	ContainerList(ctx context.Context, opts container.ListOptions) ([]container.Summary, error)
}

type DockerService struct {
	dockerClient dockerAPI
}

func NewDockerService(api dockerAPI) *DockerService {
	return &DockerService{dockerClient: api}
}

func (d *DockerService) CheckDockerDaemon(ctx context.Context) error {
	_, err := d.dockerClient.Ping(ctx)
	return err
}

func (d *DockerService) ContainersList(ctx context.Context) ([]container.Summary, error) {
	ContainersList, err := d.dockerClient.ContainerList(ctx, container.ListOptions{})
	return ContainersList, err
}
