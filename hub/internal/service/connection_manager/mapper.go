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

func toGRPCCommandRequest(requestID string, request domainHub.AgentRequest) pb.ServerCommandRequest {
	return pb.ServerCommandRequest{
		RequestId:      requestID,
		Name:           request.Name,
		TimeoutSeconds: int64(request.TimeOut),
		Args:           request.Args,
	}
}

func toAgentResponse(response *pb.AgentEvent_CommandResponse) domainHub.AgentResponse {
	return domainHub.AgentResponse{
		Success:    response.CommandResponse.Success,
		Error:      response.CommandResponse.Error,
		Output:     response.CommandResponse.Output,
		ExecTimeMS: int(response.CommandResponse.ExecTimeMs),
	}
}
