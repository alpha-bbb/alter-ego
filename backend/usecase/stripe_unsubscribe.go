package usecase

import (
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/entity"
)

type IStripeUnsubscribeUseCase interface {
	Execute(account dto.Account) (dto.SubscribeInfo, error)
}

type StripeUnsubscribeUseCase struct {
	subscribeStripeRepository repository.ISubscribeStripeRepository
	userRepository            repository.IUserRepository
	stripeDriver              payment.IStripeDriver
}

func NewStripeUnsubscribeUseCase(
	subscribeStripeRepository repository.ISubscribeStripeRepository,
	userRepository repository.IUserRepository,
	stripeDriver payment.IStripeDriver,
) IStripeUnsubscribeUseCase {
	return &StripeUnsubscribeUseCase{
		subscribeStripeRepository: subscribeStripeRepository,
		userRepository:            userRepository,
		stripeDriver:              stripeDriver,
	}
}

func (u *StripeUnsubscribeUseCase) Execute(account dto.Account) (dto.SubscribeInfo, error) {
	user, err := u.userRepository.FindByAccount(account.PlatformType.String(), account.AccountID)
	if err != nil {
		return dto.SubscribeInfo{}, err
	}
	subscribeStripe, err := u.subscribeStripeRepository.FindByUserID(user.UserID)
	if err != nil {
		return dto.SubscribeInfo{}, err
	}
	err = u.stripeDriver.CancelSubscription(subscribeStripe.SessionID)
	if err != nil {
		return dto.SubscribeInfo{}, err
	}

	err = u.subscribeStripeRepository.UpdateStatusBySessionID(
		subscribeStripe.SessionID,
		entity.SubscribeStatusNotSubscribed,
	)
	if err != nil {
		return dto.SubscribeInfo{}, err
	}

	return dto.SubscribeInfo{
		Status:      entity.SubscribeStatusNotSubscribed,
		Message:     "",
		RedirectUrl: "",
	}, nil
}
