package usecase

import (
	"context"
	"fmt"
	"time"

	"os"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/db/db_functions"
)

// setupLogger はコンソール出力を使用してロガーを初期化します。
func setupLogger() (*zap.Logger, error) {
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	logger := zap.New(consoleCore)
	return logger, nil
}

// SubmitUserChoice はユーザー選択の送信を処理します。
func SubmitUserChoice(ctx context.Context, req *backendpb.SubmitUserChoiceRequest) (*backendpb.SubmitUserChoiceResponse, error) {
	// ロガーをセットアップします。
	logger, err := setupLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	// リクエストを検証します。
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		logger.Error("Validation failed",
			zap.String("choice", req.GetChoice()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// ユーザーの選択をログに記録する。
	logger.Info("Received user choice",
		zap.String("choice", req.GetChoice()),
		zap.Time("timestamp", time.Now()),
	)

	// MongoDBに接続します。
	client, _, coll, err := db_functions.Connect()
	if err != nil {
		logger.Error("Failed to connect to MongoDB", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect MongoDB", zap.String("error", err.Error()))
		}
	}()

	// LLMのレスポンスをMongoDBに更新します。
	_, err = db_functions.UpdateConversationChoice(ctx, coll, req)
	if err != nil {
		logger.Error("Failed to update conversation choice", zap.Error(err))
		return nil, err
	}

	// レスポンスを返します。
	return &backendpb.SubmitUserChoiceResponse{
		Status: backendpb.SubmitUserChoiceResponse_STATUS_OK,
	}, nil
}
