package usecase

import (
	"github.com/alpha-bbb/alter-ego/backend/adapter/clock"
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
)

type IStripeUnsubscribeUseCase interface {
	Execute(account dto.Account) (dto.SubscribeInfo, error)
}

type StripeUnsubscribeUseCase struct {
	clock                     clock.IClock
	subscribeStripeRepository repository.ISubscribeStripeRepository
	userRepository            repository.IUserRepository
	stripeDriver              payment.IStripeDriver
}

func NewStripeUnsubscribeUseCase(
	clock clock.IClock,
	subscribeStripeRepository repository.ISubscribeStripeRepository,
	userRepository repository.IUserRepository,
	stripeDriver payment.IStripeDriver,
) IStripeUnsubscribeUseCase {
	return &StripeUnsubscribeUseCase{
		clock:                     clock,
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
	now := u.clock.Now().Unix()
	subscribeStatus, CurrentPeriodEnd, status, err := u.stripeDriver.GetSubscriptionStatus(subscribeStripe.SessionID, now)
	if err != nil {
		return dto.SubscribeInfo{}, err
	}

	err = u.subscribeStripeRepository.UpdateStatusBySessionID(
		subscribeStripe.SessionID,
		subscribeStatus,
	)
	if err != nil {
		return dto.SubscribeInfo{}, err
	}

	return dto.SubscribeInfo{
		Status:      subscribeStatus,
		Message:     status,
		RedirectUrl: "",
		ExpiresAt:   CurrentPeriodEnd,
	}, nil
}
