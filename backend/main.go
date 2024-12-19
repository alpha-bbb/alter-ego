package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/alpha-bbb/alter-ego/backend/server"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func main() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    grpcClientAddress := os.Getenv("LLM_GRPC_CLIENT_ADDRESS")
    port := os.Getenv("PORT")

    // Setup dependencies
    llmClient := server.NewGRPCLLMClient(grpcClientAddress)
    talkUseCase := usecase.NewTalkInteractor(llmClient)
    backendServer := server.NewBackendServer(talkUseCase)

    // Start gRPC server
    listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    backendpb.RegisterBackendServiceServer(s, backendServer)
    reflection.Register(s)

    go func() {
        log.Printf("starting gRPC server on port %s", port)
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
