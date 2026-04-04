package main

import (
	"log"
	"net"

	"github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	grpcserver "github.com/lorsanstand/HomeOps-Hub/internal/hub/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":6756")
	if err != nil {
		return
	}

	grpcServer := grpc.NewServer()

	srv := &grpcserver.Server{}
	homeops.RegisterHubServer(grpcServer, srv)

	log.Println("Start serve")
	grpcServer.Serve(lis)
}
