package rpc

import (
	"context"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/shared/domain"
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
	ResponseData, err := c.Hub().RegisterAgent(ctx, new(domain.ToGRPCAgentRequest(RegisterData)))
	return domain.ToDomainAgentResponse(ResponseData), err
}
