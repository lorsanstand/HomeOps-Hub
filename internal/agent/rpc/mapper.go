package rpc

import (
	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/domain"
)

func toAgentRegisterRequest(request domain.RegisterAgentData) pb.RegisterAgentRequest {
	return pb.RegisterAgentRequest{
		AgentId:   request.AgentId,
		AgentName: request.AgentName,
		Host: &pb.HostInfo{
			Hostname: request.Host.Hostname,
			Arch:     request.Host.Arch,
			System:   request.Host.System,
		},
		Version:    request.AgentVersion,
		Capability: toGRPCCapability(request.Capabilities),
	}
}

func toGRPCCapability(caps []domain.Capability) []*pb.Capability {
	var capability []*pb.Capability
	for _, capi := range caps {
		capability = append(capability, &pb.Capability{
			Name:      capi.Name,
			Available: capi.Available,
			Version:   capi.Version,
			Reason:    capi.Reason,
		})
	}
	return capability
}

func toAgentRegisterDataResponse(response *pb.RegisterAgentResponse) domain.RegisterAgentDataResponse {
	if response == nil {
		return domain.RegisterAgentDataResponse{}
	}

	return domain.RegisterAgentDataResponse{
		AgentID:   response.AgentId,
		Heartbeat: int(response.HeartbeatIntervalSecond),
	}
}
