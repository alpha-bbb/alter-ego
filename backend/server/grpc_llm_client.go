package server

import (
	"context"

	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type LLMClient interface {
	Talk(ctx context.Context, req *llmpb.TalkRequest) (*llmpb.TalkResponse, error)
}

type grpcLLMClient struct {
	client llmpb.LlmServiceClient
}

func NewGRPCLLMClient(address string) LLMClient {
	logger, err := log.NewLogger()
	if err != nil {
		logger.Fatal("failed to create logger", zap.Error(err))
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("failed to connect to LLM server", zap.Error(err))
	}
	return &grpcLLMClient{client: llmpb.NewLlmServiceClient(conn)}
}

func (c *grpcLLMClient) Talk(ctx context.Context, req *llmpb.TalkRequest) (*llmpb.TalkResponse, error) {
	return c.client.Talk(ctx, req)
}
