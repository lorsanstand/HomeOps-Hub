package rpc

import (
	"context"
	"fmt"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
	"github.com/lorsanstand/HomeOps-Hub/shared/mappers/rpc"
	"google.golang.org/grpc"
)

type Connection struct {
	hub  pb.HubClient
	conn *grpc.ClientConn
}

func NewConnectAgent(conn *grpc.ClientConn) *Connection {
	client := pb.NewHubClient(conn)
	return &Connection{hub: client, conn: conn}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) Hub() pb.HubClient {
	return c.hub
}

func (c *Connection) RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentRequest) (domain.RegisterAgentResponse, error) {
	ResponseData, err := c.Hub().RegisterAgent(ctx, rpc.ToGRPCAgentRequest(RegisterData))
	if err != nil {
		return domain.RegisterAgentResponse{}, fmt.Errorf("send register agent: %w", err)
	}
	response, err := rpc.ToDomainAgentResponse(ResponseData)
	if err != nil {
		return domain.RegisterAgentResponse{}, fmt.Errorf("casting response: %w", err)
	}
	return response, nil
}
