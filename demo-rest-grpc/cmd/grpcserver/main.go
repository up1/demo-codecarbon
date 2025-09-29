package main

import (
	"context"
	"log"
	"net"

	hello "demo/hellobench/gen/hello"

	"google.golang.org/grpc"
)

type server struct {
	hello.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloReply, error) {
	msg := req.GetMessage()
	if msg == "" {
		msg = "hello world"
	}
	return &hello.HelloReply{Message: msg}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	hello.RegisterHelloServiceServer(grpcServer, &server{})
	log.Printf("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
