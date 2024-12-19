package server

import (
	"context"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
)

type TalkUseCase interface {
	HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error)
}

type BackendServer struct {
	backendpb.UnimplementedBackendServiceServer
	talkUseCase TalkUseCase
}

func NewBackendServer(talkUseCase TalkUseCase) *BackendServer {
	return &BackendServer{talkUseCase: talkUseCase}
}

func (s *BackendServer) Talk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
	return s.talkUseCase.HandleTalk(ctx, req)
}
