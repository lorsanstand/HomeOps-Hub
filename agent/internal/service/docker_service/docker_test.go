package docker_service

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog"
)

var errTest error = errors.New("test")

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
	t.Parallel()

	tests := []struct {
		name    string
		mock    DockerMock
		wantErr error
	}{
		{
			name: "success",
			mock: DockerMock{
				pingErr:      nil,
				containers:   nil,
				containerErr: nil,
			},
			wantErr: nil,
		},
		{
			name: "docker error",
			mock: DockerMock{
				pingErr:      errTest,
				containers:   nil,
				containerErr: nil,
			},
			wantErr: errTest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewDockerService(tt.mock, zerolog.Logger{})

			err := svc.CheckDockerDaemon(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got: %v", tt.wantErr, err)
			}
		})
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
				containerErr: errTest,
			},
			wantLen: 0,
			wantErr: errTest,
		},
		{
			name: "docker empty container",
			mock: DockerMock{
				pingErr:      nil,
				containers:   nil,
				containerErr: nil,
			},
			wantLen: 0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewDockerService(tt.mock, zerolog.Logger{})

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
