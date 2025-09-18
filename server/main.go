package main

import (
	"log"
	"net"

	accountv1 "github.com/galadeat/bank-sim/api/proto/account/v1"
	"github.com/galadeat/bank-sim/server/internal"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	accountv1.RegisterAccountServer(s, internal.NewServer())

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
