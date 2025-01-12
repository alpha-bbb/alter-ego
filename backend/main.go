package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/alpha-bbb/alter-ego/backend/server"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func main() {
	logger, err := log.NewLogger()
	if err != nil {
		logger.Fatal("failed to create logger: %v", zap.Error(err))
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
		logger.Fatal("failed to listen: %v", zap.Error(err))
	}
	s := grpc.NewServer()
	backendpb.RegisterBackendServiceServer(s, backendServer)
	reflection.Register(s)

	go func() {
		logger.Info("starting gRPC server on port", zap.String("port", port))
		if err := s.Serve(listener); err != nil {
			logger.Fatal("failed to serve: %v", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("stopping gRPC server...")
	s.GracefulStop()
}
