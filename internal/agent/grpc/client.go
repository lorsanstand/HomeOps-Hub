package grpc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HomeOpsAgent struct {
	conn pb.HubClient
}

func NewConnectAgent(address string) (*HomeOpsAgent, error) {
	conn, err := grpc.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed connection hub: %v", err)
	}

	client := pb.NewHubClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed connection hub: %v", err)
	}

	if resp.Pong != "Pong" {
		return nil, fmt.Errorf("failed connection hub: %v", err)
	}

	return &HomeOpsAgent{conn: client}, nil
}
