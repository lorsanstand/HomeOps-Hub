package rpc

import (
	"context"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"github.com/lorsanstand/HomeOps-Hub/internal/agent/domain"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Connection struct {
	hub  pb.HubClient
	conn *grpc.ClientConn
	log  zerolog.Logger
}

func NewConnectAgent(conn *grpc.ClientConn, logger zerolog.Logger) *Connection {
	logger = logger.With().Str("component", "agent.rpc").Logger()

	client := pb.NewHubClient(conn)

	return &Connection{hub: client, conn: conn, log: logger}
}

func (c *Connection) Close() error {
	c.log.Warn().Msg("connection close")
	return c.conn.Close()
}

func (c *Connection) Hub() pb.HubClient {
	return c.hub
}

func (c *Connection) RegisterAgent(ctx context.Context, RegisterData domain.RegisterAgentData) (domain.RegisterAgentDataResponse, error) {
	ResponseData, err := c.Hub().RegisterAgent(ctx, new(toAgentRegisterRequest(RegisterData)))
	c.log.Info().Msg("register agent")
	return toAgentRegisterDataResponse(ResponseData), err
}
