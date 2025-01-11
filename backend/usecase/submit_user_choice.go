package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

func SubmitUserChoice(ctx context.Context, req *backendpb.SubmitUserChoiceRequest) (*backendpb.SubmitUserChoiceResponse, error) {
    // ログファイルの作成またはオープン
    file, err := os.OpenFile("/var/log/alter-ego/user_choices.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

    if err != nil {
        fmt.Printf("Failed to open log file: %v\n", err)
        return nil, fmt.Errorf("failed to open log file: %w", err)
    }
    defer file.Close()

    // ファイルハンドラーを用いたzap loggerの設定
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
        zapcore.AddSync(file),
        zap.InfoLevel,
    )
    logger := zap.New(core)

    defer func() {
        if err := logger.Sync(); err != nil {
            fmt.Printf("Failed to sync logger: %v\n", err)
        }
    }()

    // バリデーション処理
    validate := validator.New()
    err = validate.Struct(req)
    if err != nil {
        logger.Error("validation failed",
            zap.String("choice", req.GetChoice()),
            zap.Error(err),
        )
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 受け取った選択肢をログに出力
    logger.Info("Received user choice",
        zap.String("choice", req.GetChoice()),
        zap.Time("timestamp", time.Now()),
    )

    return &backendpb.SubmitUserChoiceResponse{
        Status: backendpb.SubmitUserChoiceResponse_STATUS_OK,
    }, nil
}
