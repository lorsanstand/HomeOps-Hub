package main

import (
	"context"
	"log"
	"time"

	"github.com/lorsanstand/HomeOps-Hub/api/gen/homeops"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:6756", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := homeops.NewHubClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("dial: %v", err)
	}

	defer cancel()

	log.Printf("pong: %+v", resp.Pong)
}
