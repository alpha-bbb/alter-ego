package usecase

import (
	"context"
	"fmt"
	"strings"

	convert "github.com/alpha-bbb/alter-ego/backend/convert"
	"github.com/alpha-bbb/alter-ego/backend/preprocess"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"github.com/alpha-bbb/alter-ego/backend/server"
)

type TalkUseCase interface {
	HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error)
}

type talkInteractor struct {
	llmClient server.LLMClient
}

func NewTalkInteractor(llmClient server.LLMClient) TalkUseCase {
	return &talkInteractor{llmClient: llmClient}
}

func (i *talkInteractor) HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
	// 1. リクエストをentityTalkHistoryに変換
	entityTalkHistory := convert.ConvertTalkHistoryFromGRPCTalkRequest(req)

	// 2. Messageフィールドをトークン化
	for _, history := range entityTalkHistory {
		history.Message = preprocess.TokenizeMessage(history.Message)
	}

	// 3. トークン化されたデータをGRPC形式に変換
	llmHistories := convert.ConvertTalkHistoryToGRPCTalkResponse(entityTalkHistory)

	// 4. LLMサービスにリクエスト
	llmRequest := &llmpb.TalkRequest{Histories: llmHistories}
	llmResponse, err := i.llmClient.Talk(ctx, llmRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM service: %w", err)
	}

	// 5. レスポンスを復元
	restoredMessage := preprocess.RestoreTokens(strings.Join(llmResponse.Message, " "))

	// 6. トークンをクリア
	preprocess.ClearAllTokens()

	// 7. 復元したメッセージをレスポンスとして返す
	return &backendpb.TalkResponse{Message: []string{restoredMessage}}, nil
}
