package usecase

import (
	"context"
	"fmt"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"github.com/alpha-bbb/alter-ego/backend/server"
)

type TalkUseCase interface {
	HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error)
}

func convertTalkHistoryFromGRPCTalkRequest(req *backendpb.TalkRequest) []*llmpb.TalkHistory {
	// Implement the conversion logic here
	// This is a placeholder implementation
	return []*llmpb.TalkHistory{}
}

type talkInteractor struct {
	llmClient server.LLMClient
}
func NewTalkInteractor(llmClient server.LLMClient) TalkUseCase {
	return &talkInteractor{llmClient: llmClient}
}

func (i *talkInteractor) HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
	entityTalkHistory := convertTalkHistoryFromGRPCTalkRequest(req)
	llmRequest := &llmpb.TalkRequest{Histories: entityTalkHistory}
	llmResponse, err := i.llmClient.Talk(ctx, llmRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM service: %w", err)
	}
	return &backendpb.TalkResponse{Message: llmResponse.Message}, nil
}
