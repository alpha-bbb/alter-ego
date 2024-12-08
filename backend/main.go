package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type BackendServer struct {
	backendpb.UnimplementedBackendServiceServer
}

func ConvertTalkHistoryFromGRPCTalkRequest(req *backendpb.TalkRequest) []*llmpb.TalkHistory {
	res := make([]*llmpb.TalkHistory, len(req.Histories))
	for i := 0; i < len(req.Histories); i++ {
		res[i] = &llmpb.TalkHistory{
			Date: req.Histories[i].Date,
			User: &llmpb.User{
				UserId: req.Histories[i].User.UserId,
				Name:   req.Histories[i].User.Name,
			},
			Message: req.Histories[i].Message,
		}
	}
	return res
}
func (s *BackendServer) Talk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
    // gRPC サーバーのアドレス
    const serverAddress = "localhost:50051"

    // gRPC クライアントを作成
    client, cleanup, err := newGRPCClient(serverAddress)
    if err != nil {
        log.Fatalf("failed to create gRPC client: %v", err)
    }
    defer cleanup()

    // LlmServiceのTalkメソッドを呼び出す
    llmResponse, err := callCheck(client, req)
    if err != nil {
        return nil, fmt.Errorf("failed to call llmService.Talk: %v", err)
    }

    // BackendServiceのレスポンスとしてllmResponseからメッセージを設定
    return &backendpb.TalkResponse{
		Message: llmResponse.Message,
    }, nil
}

// newGRPCClient は、新しい gRPC クライアントを作成します
func newGRPCClient(serverAddress string) (llmpb.LlmServiceClient, func(), error) {
    // サーバーへ接続
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, err
    }

    // クリーンアップ関数を定義
    cleanup := func() {
        if err := conn.Close(); err != nil {
            log.Printf("failed to close connection: %v", err)
        }
    }

    // LlmServiceClient を作成して返却
    return llmpb.NewLlmServiceClient(conn), cleanup, nil
}

// Check メソッドの実行
// Check メソッドの実行
func callCheck(client llmpb.LlmServiceClient, req *backendpb.TalkRequest) (*llmpb.TalkResponse, error) {
    // タイムアウトを設定
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    llmRequest := &llmpb.TalkRequest{
        Histories: ConvertTalkHistoryFromGRPCTalkRequest(req),
    }

    res, err := client.Talk(ctx, llmRequest)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func NewBackendServer() *BackendServer {
	return &BackendServer{}
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 3. gRPCサーバーにGreetingServiceを登録
	backendpb.RegisterBackendServiceServer(s, NewBackendServer())

	// 4. サーバーリフレクションの設定
	reflection.Register(s)

	// 5. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 6.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
