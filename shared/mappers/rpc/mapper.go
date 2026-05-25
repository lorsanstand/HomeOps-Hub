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

func ToDomainCapabilities(capabilities []*pb.Capability) []domain.Capability {
	domainCaps := make([]domain.Capability, len(capabilities))

	for id, capability := range capabilities {
		if capability == nil {
			continue
		}

		domainCaps[id] = domain.Capability{
			Name:      capability.Name,
			Version:   capability.Version,
			Reason:    capability.Reason,
			Available: capability.Available,
			Command:   ToDomainCapabilityCommands(capability.Command),
		}
	}

	return domainCaps
}

func ToDomainCapabilityCommands(commands []*pb.CapabilityCommand) []domain.CapabilityCommand {
	domainCommand := make([]domain.CapabilityCommand, len(commands))

	for id, command := range commands {
		if command == nil {
			continue
		}

		domainCommand[id] = domain.CapabilityCommand{
			Name:         command.Name,
			OptionalArgs: ToDomainCommandArgs(command.OptArgs),
			RequiredArgs: ToDomainCommandArgs(command.ReqArgs),
			Version:      command.Version,
		}
	}
	return domainCommand
}

func ToDomainCommandArgs(args []*pb.CommandsArgs) []domain.CommandArgs {
	DomainArgs := make([]domain.CommandArgs, len(args))

	for id, arg := range args {
		DomainArgs[id] = domain.CommandArgs{
			Name:        arg.Name,
			Type:        arg.Type,
			Default:     arg.Default,
			Description: arg.Description,
			Enum:        arg.Enum,
			Validation: domain.ArgValidation{
				AllowedExts: arg.Validation.AllowedExts,
				MaxValue:    int(arg.Validation.MaxValue),
				MinValue:    int(arg.Validation.MinValue),
				Pattern:     arg.Validation.Pattern,
			},
		}
	}

	return DomainArgs
}

func ToGRPCAgentRequest(request domain.RegisterAgentRequest) *pb.RegisterAgentRequest {
	return &pb.RegisterAgentRequest{
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

func ToGRPCCapability(capabilities []domain.Capability) []*pb.Capability {
	GRPCCapabilities := make([]*pb.Capability, len(capabilities))

	for id, capability := range capabilities {
		GRPCCapabilities[id] = &pb.Capability{
			Name:      capability.Name,
			Available: capability.Available,
			Version:   capability.Version,
			Reason:    capability.Reason,
			Command:   ToGRPCCapabilityCommands(capability.Command),
		}
	}
	return GRPCCapabilities
}

func ToGRPCCapabilityCommands(commands []domain.CapabilityCommand) []*pb.CapabilityCommand {
	GRPCCommands := make([]*pb.CapabilityCommand, len(commands))

	for id, command := range commands {
		GRPCCommands[id] = &pb.CapabilityCommand{
			Name:       command.Name,
			Version:    command.Version,
			OptArgs:    ToGRPCCommandArgs(command.OptionalArgs),
			ReqArgs:    ToGRPCCommandArgs(command.RequiredArgs),
			TypeOutput: command.TypeOutput,
		}
	}

	return GRPCCommands
}

func ToGRPCCommandArgs(args []domain.CommandArgs) []*pb.CommandsArgs {
	GRPCArgs := make([]*pb.CommandsArgs, len(args))

	for id, arg := range args {
		GRPCArgs[id] = &pb.CommandsArgs{
			Name:        arg.Name,
			Type:        arg.Type,
			Default:     arg.Default,
			Description: arg.Description,
			Enum:        arg.Enum,
			Validation: &pb.ArgValidation{
				AllowedExts: arg.Validation.AllowedExts,
				MaxValue:    int64(arg.Validation.MaxValue),
				MinValue:    int64(arg.Validation.MinValue),
				Pattern:     arg.Validation.Pattern,
			},
		}
	}

	return GRPCArgs
}
