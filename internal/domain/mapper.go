package domain

import (
	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
)

func ToDomainAgentRequest(request *pb.RegisterAgentRequest) RegisterAgentData {
	if request == nil {
		return RegisterAgentData{}
	}

	return RegisterAgentData{
		AgentId:   request.AgentId,
		AgentName: request.AgentName,
		Host: HostInfo{
			System:   request.Host.System,
			Hostname: request.Host.Hostname,
			Arch:     request.Host.Arch,
		},
		Capabilities: ToDomainCapabilities(request.Capability),
	}
}

func ToDomainAgentResponse(response *pb.RegisterAgentResponse) RegisterAgentDataResponse {
	if response == nil {
		return RegisterAgentDataResponse{}
	}

	return RegisterAgentDataResponse{
		AgentID:   response.AgentId,
		Heartbeat: int(response.HeartbeatIntervalSecond),
	}
}

func ToDomainCapabilities(capability []*pb.Capability) []Capability {
	var caps []Capability

	for _, capa := range capability {
		if capa == nil {
			continue
		}

		caps = append(caps, Capability{
			Name:      capa.Name,
			Version:   capa.Version,
			Reason:    capa.Reason,
			Available: capa.Available,
		})
	}

	return caps
}

func ToGRPCAgentRequest(request RegisterAgentData) pb.RegisterAgentRequest {
	return pb.RegisterAgentRequest{
		AgentId:   request.AgentId,
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

func ToGRPCCapability(caps []Capability) []*pb.Capability {
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
