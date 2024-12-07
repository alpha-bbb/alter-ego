package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type BackendServer struct {
	backendpb.UnimplementedBackendServiceServer
}

func ConvertTalkHistoryFromGRPCTalkRequest(req *backendpb.TalkRequest) []entity.TalkHistory {
	res := make([]entity.TalkHistory, len(req.Histories))
	for i := 0; i < len(req.Histories); i++ {
        res[i] = entity.TalkHistory{
			Date:    req.Histories[i].Date,
			User:    entity.User{
				UserID:   req.Histories[i].User.UserId,
				Name: req.Histories[i].User.Name,
			},
			Message: req.Histories[i].Message,
		}
    }
	return res
}

func (s *BackendServer) Talk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
	fmt.Printf("%w",(ConvertTalkHistoryFromGRPCTalkRequest(req)))
	return &backendpb.TalkResponse{
		Message: []string{"Hello"},
	}, nil
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