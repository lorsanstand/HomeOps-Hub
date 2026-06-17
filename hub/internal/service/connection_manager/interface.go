package connection_manager

import (
	"context"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

type StreamConn interface {
	Send(request *pb.ServerCommandRequest) error
	Recv() (*pb.AgentEvent, error)
	Context() context.Context
}

type heartbeatStore interface {
	CreateHeartbeat(ctx context.Context, heartbeat domainHub.CreateHeartbeatModel) error
}

type StatusAgent interface {
	Offline()
	Online()
}

type statusNotifier interface {
	New(agentID string) StatusAgent
}
