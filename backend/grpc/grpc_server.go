package grpc

import (
	"fmt"
	"net"

	"github.com/alpha-bbb/alter-ego/backend/config"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/client"
	"github.com/alpha-bbb/alter-ego/backend/usecase"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(logger *zap.Logger) (*grpc.Server, net.Listener, error) {
	grpcPort := config.GrpcPort()

	grpcClientAddress := config.LlmGrpcClientAddress()
	llmClient := client.NewGRPCLLMClient(grpcClientAddress)
	talkUseCase := usecase.NewTalkUseCase(llmClient)
	backendServer := NewBackendServiceServer(talkUseCase)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()
	backendpb.RegisterBackendServiceServer(grpcServer, backendServer)
	reflection.Register(grpcServer)

	return grpcServer, lis, nil
}
