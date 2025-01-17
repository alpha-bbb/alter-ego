package usecase

import (
	"context"
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/convert"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	db_convert "github.com/alpha-bbb/alter-ego/backend/infrastructure/db/db_convert"
	db_functions "github.com/alpha-bbb/alter-ego/backend/infrastructure/db/db_functions"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/alpha-bbb/alter-ego/backend/server"
	"go.uber.org/zap"
)

func HandleTalk(ctx context.Context, req *backendpb.TalkRequest, llmClient server.LLMClient) (*backendpb.TalkResponse, error) {
	// ロガーの生成
	logger, err := log.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	logger.Info("Received talk request")

	// gRPC リクエストから内部の会話履歴エンティティに変換
	entityTalkHistory := convert.ConvertTalkHistoryFromGRPCTalkRequest(req)

	// MongoDB に接続
	client, _, coll, err := db_functions.Connect()
	if err != nil {
		logger.Info(err.Error())
		logger.Error("Failed to connect to MongoDB")
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect MongoDB", zap.String("error", err.Error()))
		}
	}()

	// 会話履歴エンティティをモデルに変換
	newConv := db_convert.ConvertEntityToModel(entityTalkHistory, logger)

	// 会話履歴の挿入
	insertedIDInterface, err := db_functions.InsertConversation(ctx, coll, *newConv)
	insertedID := insertedIDInterface.(string)
	if err != nil {
		logger.Error("Failed to insert conversation", zap.Error(err))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Inserted conversation: %s", insertedID))

	// LLM への変換
	llmHistories := convert.ConvertTalkHistoryToGRPCTalkResponse(entityTalkHistory)
	llmRequest := &llmpb.TalkRequest{Histories: llmHistories}

	// LLM へのリクエスト
	llmResponse, err := llmClient.Talk(ctx, llmRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM service: %w", err)
	}

	// LLM のレスポンスを MongoDB に更新
	_, err = db_functions.UpdateLlmResponse(ctx, coll, newConv.ID, llmResponse.Message)
	if err != nil {
		logger.Error("Failed to update conversation choice", zap.Error(err))
		return nil, err
	}

	return &backendpb.TalkResponse{ConversationId: insertedID, Message: llmResponse.Message}, nil
}
