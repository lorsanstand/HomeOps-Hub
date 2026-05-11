package connection_manager

import (
	"time"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	domainHub "github.com/lorsanstand/HomeOps-Hub/hub/internal/domain"
)

func toCreateHeartbeatModel(agentID string, heartbeat *pb.AgentEvent_Heartbeat) domainHub.CreateHeartbeatModel {
	timestamp := time.Unix(heartbeat.Heartbeat.Timestamp, 0)

	return domainHub.CreateHeartbeatModel{
		AgentID:   agentID,
		Timestamp: timestamp,
		Metrics: domainHub.SystemMetrics{
			MemoryUsage: float64(heartbeat.Heartbeat.Metrics.MemoryUsage),
			CpuUsage:    float64(heartbeat.Heartbeat.Metrics.CpuUsage),
			DiskUsage:   float64(heartbeat.Heartbeat.Metrics.DiskUsage),
		},
	}
}
