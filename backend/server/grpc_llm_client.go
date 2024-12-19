package server

import (
	"context"
	"log"

	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"google.golang.org/grpc"
)

type LLMClient interface {
	Talk(ctx context.Context, req *llmpb.TalkRequest) (*llmpb.TalkResponse, error)
}

type grpcLLMClient struct {
	client llmpb.LlmServiceClient
}

func NewGRPCLLMClient(address string) LLMClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to LLM server: %v", err)
	}
	return &grpcLLMClient{client: llmpb.NewLlmServiceClient(conn)}
}

func (c *grpcLLMClient) Talk(ctx context.Context, req *llmpb.TalkRequest) (*llmpb.TalkResponse, error) {
	return c.client.Talk(ctx, req)
}
