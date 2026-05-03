package rpc

import (
	"context"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HubService interface {
	RegisterAgent(ctx context.Context, data domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error)
}

type HubHandler struct {
	pb.UnimplementedHubServer
	log        zerolog.Logger
	GrpcServer *grpc.Server
	hub        HubService
}

func NewHubHandler(HubServ HubService, logger zerolog.Logger) *HubHandler {
	hub := &HubHandler{log: logger, hub: HubServ}

	grpcServer := grpc.NewServer()
	pb.RegisterHubServer(grpcServer, hub)

	hub.GrpcServer = grpcServer

	return hub
}

func (h *HubHandler) Ping(ctx context.Context, _ *emptypb.Empty) (*pb.PongResponse, error) {
	h.log.Debug().Msg("ping request received")
	return &pb.PongResponse{Pong: "Pong"}, nil
}

func (h *HubHandler) RegisterAgent(ctx context.Context, request *pb.RegisterAgentRequest) (*pb.RegisterAgentResponse, error) {
	h.log.Debug().Str("agentId", request.AgentId).Str("agentName", request.AgentName).Msg("register agent request received")
	data := domain.ToDomainAgentRequest(request)
	resp, err := h.hub.RegisterAgent(ctx, data)
	if err != nil {
		h.log.Error().Err(err).Str("agentId", request.AgentId).Msg("register agent request failed")
		return domain.ToGRPCAgentResponse(resp), err
	}
	h.log.Debug().Str("agentId", resp.AgentID).Msg("register agent request completed")
	return domain.ToGRPCAgentResponse(resp), nil
}
