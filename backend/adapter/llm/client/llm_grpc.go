package client

import (
	"context"

	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ILLMClient interface {
	Talk(req *llmpb.TalkRequest) (*llmpb.TalkResponse, error)
}

type grpcLLMClient struct {
	client llmpb.LlmServiceClient
}

func NewLLMGRPCClient(address string) ILLMClient {
	logger, err := log.NewLogger()
	if err != nil {
		logger.Fatal("failed to create logger", zap.Error(err))
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to connect to LLM server", zap.Error(err))
	}
	return &grpcLLMClient{
		client: llmpb.NewLlmServiceClient(conn),
	}
}

func (c *grpcLLMClient) Talk(req *llmpb.TalkRequest) (*llmpb.TalkResponse, error) {
	return c.client.Talk(context.Background(), req)
}
