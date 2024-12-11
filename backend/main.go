package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/alpha-bbb/alter-ego/backend/server"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func main() {
	// Setup dependencies
	llmClient := server.NewGRPCLLMClient("localhost:8080")
	talkUseCase := usecase.NewTalkInteractor(llmClient)
	backendServer := server.NewBackendServer(talkUseCase)

	// Start gRPC server
	port := 50051
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	backendpb.RegisterBackendServiceServer(s, backendServer)
	reflection.Register(s)

	go func() {
		log.Printf("starting gRPC server on port %d", port)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
