package hub_service

import (
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/rpc"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/service/docker_service"
	"github.com/rs/zerolog"
)

type HubService struct {
	docker  *docker_service.DockerService
	log     zerolog.Logger
	hubConn *rpc.Connection
}

func NewHubService(docker *docker_service.DockerService, log zerolog.Logger) *HubService {
	return &HubService{docker: docker, log: log}
}

func (h *HubService) GatherInfoSystem() {

}
