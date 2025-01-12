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
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
)

const (
	logDir      = "./logs"
	logFilePath = "./logs/user_choices.log"
)

// createLogDirectory ensures the log directory exists.
func createLogDirectory() error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}
	}
	return nil
}

// createLogFile ensures the log file exists.
func createLogFile() error {
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		file, err := os.Create(logFilePath)
		if err != nil {
			return fmt.Errorf("failed to create log file: %w", err)
		}
		file.Close()
	}
	return nil
}

// setupLogger initializes the logger with file and console output.
func setupLogger() (*zap.Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		zap.InfoLevel,
	)

	logger := zap.New(zapcore.NewTee(fileCore, consoleCore))
	return logger, nil
}

// SubmitUserChoice handles user choice submissions.
func SubmitUserChoice(ctx context.Context, req *backendpb.SubmitUserChoiceRequest) (*backendpb.SubmitUserChoiceResponse, error) {
	logger, err := log.NewLogger()
	if err != nil {
		return nil, err
	}
	// Ensure log directory and file exist.
	if err := createLogDirectory(); err != nil {
		logger.Error("Failed to create log directory: %v\n", zap.Error(err))
		return nil, err
	}

	if err := createLogFile(); err != nil {
		logger.Error("Failed to create log file: %v\n", zap.Error(err))
		return nil, err
	}

	// Set up the logger.
	logger_store, err := setupLogger()
	if err != nil {
		logger.Error("Failed to setup logger: %v\n", zap.Error(err))
		return nil, err
	}

	// Validate the request.
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		logger_store.Error("Validation failed",
			zap.String("choice", req.GetChoice()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Log the user choice.
	logger_store.Info("Received user choice",
		zap.String("choice", req.GetChoice()),
		zap.Time("timestamp", time.Now()),
	)

	// Return the response.
	return &backendpb.SubmitUserChoiceResponse{
		Status: backendpb.SubmitUserChoiceResponse_STATUS_OK,
	}, nil
}
