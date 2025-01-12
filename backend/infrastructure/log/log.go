package log

import (
	"os"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func NewLogger() (*zap.Logger, error) {
	if logger != nil {
		return logger, nil
	}

	var err error

	env := os.Getenv("ENV")

	switch env {
	case "production":
		logger, err = zap.NewProduction()
	case "development":
		logger, err = zap.NewDevelopment()
	case "test":
		logger = zap.NewNop()
	default:
		logger, err = zap.NewDevelopment()
	}

	logger.Info(env)

	return logger, err
}
