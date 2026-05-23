package rpc

import (
	"context"
	"fmt"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/hub/internal/service/connection_manager"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HubService interface {
	RegisterAgent(ctx context.Context, data domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error)
}

type ConnectionManager interface {
	NewConnection(stream connection_manager.StreamConn) error
}

type HubHandler struct {
	pb.UnimplementedHubServer
	log           zerolog.Logger
	GrpcServer    *grpc.Server
	hub           HubService
	streamManager ConnectionManager
}

func NewHubHandler(HubServ HubService, manager ConnectionManager, logger zerolog.Logger) *HubHandler {
	hub := &HubHandler{log: logger, hub: HubServ, streamManager: manager}

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
	h.log.Debug().Str("agentID", request.AgentId).Str("agentName", request.AgentName).Msg("register agent request received")
	data, err := domain.ToDomainAgentRequest(request)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to casting request")
		return &pb.RegisterAgentResponse{}, ErrFailedRegister
	}
	resp, err := h.hub.RegisterAgent(ctx, data)
	if err != nil {
		h.log.Error().Err(err).Msg("failed register agent")
		return domain.ToGRPCAgentResponse(resp), ErrFailedRegister
	}
	h.log.Info().Str("agentID", resp.AgentID).Msg("register agent request completed")
	return domain.ToGRPCAgentResponse(resp), nil
}

func (h *HubHandler) StreamConnection(stream grpc.BidiStreamingServer[pb.AgentEvent, pb.ServerCommandRequest]) error {
	err := h.streamManager.NewConnection(stream)
	if err != nil {
		return fmt.Errorf("accept connection: %w", err)
	}
	return nil
}
