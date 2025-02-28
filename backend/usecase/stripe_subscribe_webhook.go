package usecase

import (
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
)

type IStripeSubscribeWebhookUseCase interface {
	Execute(dto dto.StripeSubscribeWebhook) error
}

type StripeSubscribeWebhookUseCase struct {
	subscribeStripeRepository repository.ISubscribeStripeRepository
	stripeDriver              payment.IStripeDriver
}

func NewStripeSubscribeWebhookUseCase(
	subscribeStripeRepository repository.ISubscribeStripeRepository,
	stripeDriver payment.IStripeDriver,
) IStripeSubscribeWebhookUseCase {
	return &StripeSubscribeWebhookUseCase{
		subscribeStripeRepository: subscribeStripeRepository,
		stripeDriver:              stripeDriver,
	}
}

func (u *StripeSubscribeWebhookUseCase) Execute(req dto.StripeSubscribeWebhook) error {
	return u.subscribeStripeRepository.UpdateStatusBySessionID(req.SessionID, req.Status)
}
