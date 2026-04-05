package rpc

import (
	"fmt"

	pb "github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"google.golang.org/grpc"
)

type Connection struct {
	hub  pb.HubClient
	conn *grpc.ClientConn
}

func NewConnectAgent(address string) (*Connection, error) {
	conn, err := grpc.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed connection hub: %v", err)
	}

	client := pb.NewHubClient(conn)

	return &Connection{hub: client, conn: conn}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

func (c *Connection) Hub() pb.HubClient {
	return c.hub
}
