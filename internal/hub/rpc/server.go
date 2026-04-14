package rpc

import (
	"context"
	"fmt"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/internal/domain"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HubHandler struct {
	pb.UnimplementedHubServer
	log        zerolog.Logger
	GrpcServer *grpc.Server
}

func NewHubHandler(logger zerolog.Logger) *HubHandler {
	hub := &HubHandler{log: logger}

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
	fmt.Println(data)
	return &pb.RegisterAgentResponse{}, nil
}
