package grpc

import (
	"context"

	"github.com/alpha-bbb/alter-ego/backend/adapter/backend/converter"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
)

type TalkHandler struct {
	backendpb.UnimplementedBackendServiceServer
	talkUseCase usecase.ITalkUseCase
}

func NewBackendServiceServer(talkUseCase usecase.ITalkUseCase) *TalkHandler {
	return &TalkHandler{
		talkUseCase: talkUseCase,
	}
}

func (h *TalkHandler) Talk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
	entity := converter.ConvertTalkHistoryFromGRPCTalkRequest(req)
	res, err := h.talkUseCase.Execute(entity)
	if err != nil {
		return nil, err
	}
	return converter.ConvertMessageDTOToGRPCTalkResponse(res), nil
}

// TODO: 実装する
func (h *TalkHandler) Subscribe(ctx context.Context, req *backendpb.SubscribeRequest) (*backendpb.SubscribeResponse, error) {
	return nil, nil
}
