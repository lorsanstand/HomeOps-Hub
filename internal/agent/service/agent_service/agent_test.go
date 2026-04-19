package agent_service

import (
	"context"

	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
)

type CollectorMock struct {
	host domain.HostInfo
	caps []domain.Capability
}

func (c *CollectorMock) GatherInfoSystem() (domain.HostInfo, []domain.Capability) {
	return c.host, c.caps
}

type ConnectionMock struct {
	regAgentErr error
	regResp     domain.RegisterAgentResponse
}

func (c *ConnectionMock) RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error) {
	return c.regResp, c.regAgentErr
}
