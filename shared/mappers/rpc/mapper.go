package rpc

import (
	"fmt"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

func ToDomainAgentRequest(request *pb.RegisterAgentRequest) (domain.RegisterAgentRequest, error) {
	if request == nil {
		return domain.RegisterAgentRequest{}, fmt.Errorf("request is empty")
	}

	return domain.RegisterAgentRequest{
		AgentVersion: request.Version,
		AgentID:      request.AgentId,
		AgentName:    request.AgentName,
		Host: domain.HostInfo{
			System:   request.Host.System,
			Hostname: request.Host.Hostname,
			Arch:     request.Host.Arch,
		},
		Capabilities: ToDomainCapabilities(request.Capability),
	}, nil
}

func ToDomainAgentResponse(response *pb.RegisterAgentResponse) (domain.RegisterAgentResponse, error) {
	if response == nil {
		return domain.RegisterAgentResponse{}, fmt.Errorf("request is empty")
	}

	return domain.RegisterAgentResponse{
		AgentID:   response.AgentId,
		Heartbeat: int(response.HeartbeatIntervalSecond),
	}, nil
}

func ToDomainCapabilities(capability []*pb.Capability) []domain.Capability {
	var caps []domain.Capability

	for _, capa := range capability {
		if capa == nil {
			continue
		}

		caps = append(caps, domain.Capability{
			Name:      capa.Name,
			Version:   capa.Version,
			Reason:    capa.Reason,
			Available: capa.Available,
		})
	}

	return caps
}

func ToGRPCAgentRequest(request domain.RegisterAgentRequest) pb.RegisterAgentRequest {
	return pb.RegisterAgentRequest{
		AgentId:   request.AgentID,
		AgentName: request.AgentName,
		Host: &pb.HostInfo{
			Hostname: request.Host.Hostname,
			Arch:     request.Host.Arch,
			System:   request.Host.System,
		},
		Version:    request.AgentVersion,
		Capability: ToGRPCCapability(request.Capabilities),
	}
}

func ToGRPCAgentResponse(response domain.RegisterAgentResponse) *pb.RegisterAgentResponse {
	return &pb.RegisterAgentResponse{AgentId: response.AgentID, HeartbeatIntervalSecond: int64(response.Heartbeat)}
}

func ToGRPCCapability(caps []domain.Capability) []*pb.Capability {
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
