package main

import (
	"flag"
	"fmt"
	"grpctask/internal/service"
	pb "grpctask/pkg/api/test"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v\n", lis.Addr())
	srv := service.New()
	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, srv)

	if err := server.Serve(lis); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
