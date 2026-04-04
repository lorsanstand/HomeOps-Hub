package grpc

import (
	"context"
	"log"

	"github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	homeops.UnimplementedHubServer
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*homeops.PongResponse, error) {
	log.Println("Answer")
	return &homeops.PongResponse{Pong: "Huiiii"}, nil

}
