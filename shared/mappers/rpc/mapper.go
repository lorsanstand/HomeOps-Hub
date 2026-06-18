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

func ToDomainCapabilities(capabilities map[string]*pb.Capability) map[string]domain.Capability {
	domainCaps := make(map[string]domain.Capability, len(capabilities))

	for name, capability := range capabilities {
		if capability == nil {
			continue
		}

		domainCaps[name] = domain.Capability{
			Version:   capability.Version,
			Reason:    capability.Reason,
			Available: capability.Available,
			Command:   ToDomainCapabilityCommands(capability.Command),
		}
	}

	return domainCaps
}

func ToDomainCapabilityCommands(commands map[string]*pb.CapabilityCommand) map[string]domain.CapabilityCommand {
	domainCommand := make(map[string]domain.CapabilityCommand, len(commands))

	for name, command := range commands {
		if command == nil {
			continue
		}

		domainCommand[name] = domain.CapabilityCommand{
			OptionalArgs: ToDomainCommandArgs(command.OptArgs),
			RequiredArgs: ToDomainCommandArgs(command.ReqArgs),
			Version:      command.Version,
		}
	}
	return domainCommand
}

func ToDomainCommandArgs(args map[string]*pb.CommandsArgs) map[string]domain.CommandArgs {
	DomainArgs := make(map[string]domain.CommandArgs, len(args))

	for name, arg := range args {
		DomainArgs[name] = domain.CommandArgs{
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

func ToGRPCCapability(capabilities map[string]domain.Capability) map[string]*pb.Capability {
	GRPCCapabilities := make(map[string]*pb.Capability, len(capabilities))

	for name, capability := range capabilities {
		GRPCCapabilities[name] = &pb.Capability{
			Available: capability.Available,
			Version:   capability.Version,
			Reason:    capability.Reason,
			Command:   ToGRPCCapabilityCommands(capability.Command),
		}
	}
	return GRPCCapabilities
}

func ToGRPCCapabilityCommands(commands map[string]domain.CapabilityCommand) map[string]*pb.CapabilityCommand {
	GRPCCommands := make(map[string]*pb.CapabilityCommand, len(commands))

	for name, command := range commands {
		GRPCCommands[name] = &pb.CapabilityCommand{
			Version:    command.Version,
			OptArgs:    ToGRPCCommandArgs(command.OptionalArgs),
			ReqArgs:    ToGRPCCommandArgs(command.RequiredArgs),
			TypeOutput: command.TypeOutput,
		}
	}

	return GRPCCommands
}

func ToGRPCCommandArgs(args map[string]domain.CommandArgs) map[string]*pb.CommandsArgs {
	GRPCArgs := make(map[string]*pb.CommandsArgs, len(args))

	for name, arg := range args {
		GRPCArgs[name] = &pb.CommandsArgs{
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
