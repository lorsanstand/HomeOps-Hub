package grpc

import (
	"context"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedHubServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*pb.PongResponse, error) {
	return &pb.PongResponse{Pong: "Pong"}, nil

}
