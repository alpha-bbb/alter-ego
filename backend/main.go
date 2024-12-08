package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	llmpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/llm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type BackendServer struct {
	backendpb.UnimplementedBackendServiceServer
}

func ConvertTalkHistoryFromGRPCTalkRequest(req *backendpb.TalkRequest) []*entity.TalkHistory {
    if req == nil || req.Histories == nil {
        return nil
    }

    result := make([]*entity.TalkHistory, len(req.Histories))
    for index := range req.Histories {
        result[index] = &entity.TalkHistory{
            Date: req.Histories[index].Date,
            User: entity.User{
                UserID: req.Histories[index].User.UserId,
                Name:   req.Histories[index].User.Name,
            },
            Message: req.Histories[index].Message,
        }
    }
    return result
}

func ConvertTalkHistoryToGRPCTalkResponse(histories []*entity.TalkHistory) []*llmpb.TalkHistory {
    if histories == nil {
        return nil
    }

    result := make([]*llmpb.TalkHistory, len(histories))
    for i := range histories {
        result[i] = &llmpb.TalkHistory{
            Date: histories[i].Date,
            User: &llmpb.User{
                UserId: histories[i].User.UserID,
                Name:   histories[i].User.Name,
            },
            Message: histories[i].Message,
        }
    }
    return result
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
    llmResponse, err := callLlmService(client, req)
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
func callLlmService(client llmpb.LlmServiceClient, req *backendpb.TalkRequest) (*llmpb.TalkResponse, error) {
	// タイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	entityTalkHistory := ConvertTalkHistoryFromGRPCTalkRequest(req)
	llmTalkHistory := ConvertTalkHistoryToGRPCTalkResponse(entityTalkHistory)
	llmRequest := &llmpb.TalkRequest{
		Histories: llmTalkHistory,
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
