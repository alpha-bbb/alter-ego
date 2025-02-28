package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alpha-bbb/alter-ego/backend/config"
	"github.com/alpha-bbb/alter-ego/backend/grpc"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/alpha-bbb/alter-ego/backend/router"
	"go.uber.org/zap"
)

func main() {
	logger, err := log.NewLogger()
	if err != nil {
		fmt.Printf("failed to create logger: %v", err)
		os.Exit(1)
	}

	grpcPort := config.GrpcPort()
	httpPort := config.HttpPort()

	// gRPC サーバーのセットアップ
	grpcServer, lis, err := grpc.NewGRPCServer(logger)
	if err != nil {
		logger.Fatal("failed to setup gRPC server", zap.Error(err))
	}
	go func() {
		logger.Info("starting gRPC server", zap.String("port", grpcPort))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	// RESTful API サーバーのセットアップ
	e := router.NewRouter()
	go func() {
		logger.Info("starting RESTful API server", zap.String("port", httpPort))
		if err := e.Start(":" + httpPort); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start RESTful API server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("stopping servers...")

	grpcServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown RESTful API server", zap.Error(err))
	}
}
