package usecase

import (
	"context"

	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"github.com/alpha-bbb/alter-ego/backend/server"
)

type TalkUseCase interface {
    HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error)
    SubmitUserChoice(ctx context.Context, req *backendpb.SubmitUserChoiceRequest) (*backendpb.SubmitUserChoiceResponse, error)
}

type talkInteractor struct {
    llmClient server.LLMClient
}

func NewTalkInteractor(llmClient server.LLMClient) TalkUseCase {
    return &talkInteractor{llmClient: llmClient}
}

// HandleTalkは別ファイルで定義
func (i *talkInteractor) HandleTalk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
    return HandleTalk(ctx, req, i.llmClient)
}

// SubmitUserChoiceは別ファイルで定義
func (i *talkInteractor) SubmitUserChoice(ctx context.Context, req *backendpb.SubmitUserChoiceRequest) (*backendpb.SubmitUserChoiceResponse, error) {
    return SubmitUserChoice(ctx, req)
}
