package docker_service

import (
	"context"
	"errors"
	"testing"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
)

var testError error = errors.New("test")

type DockerMock struct {
	pingErr      error
	containers   []container.Summary
	containerErr error
}

func (d DockerMock) Ping(ctx context.Context) (types.Ping, error) {
	return types.Ping{}, d.pingErr
}

func (d DockerMock) ContainerList(ctx context.Context, _ container.ListOptions) ([]container.Summary, error) {
	return d.containers, d.containerErr
}

func TestCheckDockerDaemon(t *testing.T) {
	api := DockerMock{containerErr: nil, pingErr: nil, containers: []container.Summary{}}
	docker := NewDockerService(api)

	if err := docker.CheckDockerDaemon(context.TODO()); err != nil {
		t.Errorf("check daemon failed: %v", err)
	}
}

func TestCheckDaemonFailed(t *testing.T) {
	api := DockerMock{containerErr: nil, pingErr: testError, containers: []container.Summary{}}
	docker := NewDockerService(api)

	if err := docker.CheckDockerDaemon(context.TODO()); !errors.Is(err, testError) {
		t.Errorf("the error does not match the one originally specified: %v received: %v", testError, err)
	}
}

func TestContainersList(t *testing.T) {
	t.Parallel()

	containers := []container.Summary{
		{ID: "123", Image: "postgres:latest"},
		{ID: "456", Image: "nginx:latest"},
	}

	tests := []struct {
		name    string
		mock    DockerMock
		wantLen int
		wantErr error
	}{
		{
			name: "success",
			mock: DockerMock{
				pingErr:      nil,
				containers:   containers,
				containerErr: nil,
			},
			wantLen: len(containers),
			wantErr: nil,
		},
		{
			name: "docker error",
			mock: DockerMock{
				pingErr:      nil,
				containers:   nil,
				containerErr: testError,
			},
			wantLen: 0,
			wantErr: testError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewDockerService(tt.mock)

			got, err := svc.ContainersList(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got: %v", tt.wantErr, err)
			}

			if tt.wantLen != len(got) {
				t.Fatalf("expected %d containers, got: %d", tt.wantLen, len(got))
			}
		})
	}
}
