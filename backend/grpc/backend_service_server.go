package grpc

import (
	"context"
	"errors"

	"github.com/alpha-bbb/alter-ego/backend/adapter/backend/converter"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	backendpb "github.com/alpha-bbb/alter-ego/backend/gen/grpc/backend/v1"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
)

type TalkHandler struct {
	backendpb.UnimplementedBackendServiceServer
	talkUseCase                 usecase.ITalkUseCase
	stripeSubscribeUseCase      usecase.IStripeSubscribeUseCase
	stripeUnscribeUseCase       usecase.IStripeUnsubscribeUseCase
	stripeCheckSubscribeUseCaSe usecase.IStripeCheckSubscribeUseCase
}

func NewBackendServiceServer(
	talkUseCase usecase.ITalkUseCase,
	stripeSubscribeUseCase usecase.IStripeSubscribeUseCase,
	stripeUnscribeUseCase usecase.IStripeUnsubscribeUseCase,
	stripeCheckSubscribeUseCaSe usecase.IStripeCheckSubscribeUseCase,

) *TalkHandler {
	return &TalkHandler{
		talkUseCase:                 talkUseCase,
		stripeSubscribeUseCase:      stripeSubscribeUseCase,
		stripeUnscribeUseCase:       stripeUnscribeUseCase,
		stripeCheckSubscribeUseCaSe: stripeCheckSubscribeUseCaSe,
	}
}

func (h *TalkHandler) Talk(ctx context.Context, req *backendpb.TalkRequest) (*backendpb.TalkResponse, error) {
	logger, _ := log.NewLogger()
	var res dto.Message
	histories, err := converter.TalkHistoriesFromGRPC(req)
	if err != nil {
		logger.Error(err.Error())
		res.Status = 2
		res.Messages = make([]string, 1)
		res.Messages[0] = err.Error()
		return converter.ToGRPCTalkResponse(res), nil
	}
	res, err = h.talkUseCase.Execute(histories)
	if err != nil {
		logger.Error(err.Error())
		if errors.Is(err, usecase.ErrCountLimitExceeded) {
			res.Status = 3
			res.Messages = make([]string, 1)
			res.Messages[0] = err.Error()
			return converter.ToGRPCTalkResponse(res), nil
		}
		res.Status = 2
		res.Messages = make([]string, 1)
		res.Messages[0] = err.Error()
		return converter.ToGRPCTalkResponse(res), nil
	}
	return converter.ToGRPCTalkResponse(res), nil
}

func (h *TalkHandler) Subscribe(ctx context.Context, req *backendpb.SubscribeRequest) (*backendpb.SubscribeResponse, error) {
	logger, _ := log.NewLogger()
	var res dto.SubscribeInfo
	subEntity, err := converter.FromGRPCSubscribeRequest(req)
	if err != nil {
		logger.Error(err.Error())
		res.Message = err.Error()
		return converter.ToGRPCSubscribeResponse(
			res,
		), nil
	}

	switch subEntity.Action {
	case dto.SubscribeActionStart:
		res, err = h.stripeSubscribeUseCase.Execute(subEntity)
		if err != nil {
			logger.Error(err.Error())
			res.Message = err.Error()
		}
	case dto.SubscribeActionCancel:
		res, err = h.stripeUnscribeUseCase.Execute(subEntity.Account)
		if err != nil {
			logger.Error(err.Error())
			res.Message = err.Error()
		}
	case dto.SubscribeActionCheck:
		res, err = h.stripeCheckSubscribeUseCaSe.Execute(subEntity.Account)
		if err != nil {
			logger.Error(err.Error())
			res.Message = err.Error()
		}
	}
	return converter.ToGRPCSubscribeResponse(
		res,
	), nil
}
