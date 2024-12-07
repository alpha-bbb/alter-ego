package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/yourproject/backend/v1"
	"github.com/yourproject/llm/v1"
	"google.golang.org/grpc"
)

const (
	backendAddress  = "localhost:50052"
	mediatorAddress = "localhost:50051"
)

type mediatorServer struct {
	llm.UnimplementedLlmServiceServer
}

func (s *mediatorServer) Talk(ctx context.Context, req *llm.TalkRequest) (*llm.TalkResponse, error) {
	// コンテキストのタイムアウトを設定
	connCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// バックエンドサーバーに接続
	conn, err := grpc.DialContext(connCtx, backendAddress, grpc.WithInsecure())
	if err != nil {
		log.Printf("[ERROR] Could not connect to backend server: %v", err)
		return nil, fmt.Errorf("could not connect to backend server: %v", err)
	}
	defer conn.Close()

	// バックエンドクライアント作成
	backendClient := backend.NewBackendServiceClient(conn)

	// バックエンドサーバーにリクエストを転送
	resp, err := backendClient.Talk(ctx, req)
	if err != nil {
		log.Printf("[ERROR] Error calling backend server: %v", err)
		return nil, fmt.Errorf("error calling backend server: %v", err)
	}

	// バックエンドのレスポンスをフロントエンドに返す
	return &llm.TalkResponse{
		Message: resp.GetMessage(),
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", mediatorAddress)
	if err != nil {
		log.Fatalf("[ERROR] Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	llm.RegisterLlmServiceServer(server, &mediatorServer{})

	log.Printf("[INFO] Mediator server is running at %s", mediatorAddress)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("[ERROR] Failed to serve: %v", err)
	}
}
