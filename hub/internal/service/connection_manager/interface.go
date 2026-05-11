package connection_manager

import (
	"context"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
)

type streamConn interface {
	Send(request *pb.ServerCommandRequest) error
	Recv() (*pb.AgentEvent, error)
	Context() context.Context
	Close() error
}

type heartbeatStore interface {
	CreateHeartbeat(ctx context.Context, heartbeat domainHub.CreateHeartbeatModel) error
}
