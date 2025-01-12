package usecase

import (
	"context"
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/convert"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/alpha-bbb/alter-ego/backend/server"
	"go.uber.org/zap"
)

func HandleTalk(ctx context.Context, req *backendpb.TalkRequest, llmClient server.LLMClient) (*backendpb.TalkResponse, error) {
	logger, err := log.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	logger.Info("Received talk request", zap.String("request", req.String()))
	// 既存の処理
	entityTalkHistory := convert.ConvertTalkHistoryFromGRPCTalkRequest(req)
	llmHistories := convert.ConvertTalkHistoryToGRPCTalkResponse(entityTalkHistory)
	llmRequest := &llmpb.TalkRequest{Histories: llmHistories}
	llmResponse, err := llmClient.Talk(ctx, llmRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM service: %w", err)
	}
	return &backendpb.TalkResponse{Message: llmResponse.Message}, nil
}
