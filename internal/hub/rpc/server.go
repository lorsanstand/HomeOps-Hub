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
	RegisterAgent(data domain.RegisterAgentRequest) domain.RegisterAgentResponse
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
	h.log.Info().Msg("pong request")
	return &pb.PongResponse{Pong: "Pong"}, nil
}

func (h *HubHandler) RegisterAgent(ctx context.Context, request *pb.RegisterAgentRequest) (*pb.RegisterAgentResponse, error) {
	data := domain.ToDomainAgentRequest(request)
	resp := h.hub.RegisterAgent(data)
	return domain.ToGRPCAgentResponse(resp), nil
}
