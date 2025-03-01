package usecase

import (
	"errors"
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/adapter/clock"
	"github.com/alpha-bbb/alter-ego/backend/adapter/llm/client"
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/adapter/ulid"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/database"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/entity"
)

type IStripeSubscribeUseCase interface {
	Execute(req dto.Subscribe) (dto.SubscribeInfo, error)
}

type StripeSubscribeUseCase struct {
	ulid                      ulid.IULID
	clock                     clock.IClock
	transaction               database.ITransaction
	stateCountRepository      repository.IStateCountRepository
	subscribeStripeRepository repository.ISubscribeStripeRepository
	userAccountLineRepository repository.IUserAccountLineRepository
	userAggregateRepository   repository.IUserAggregateRepository
	userDetailRepository      repository.IUserDetailRepository
	userRepository            repository.IUserRepository
	llmClient                 client.ILLMClient
	stripeDriver              payment.IStripeDriver
}

func NewStripeSubscribeUseCase(
	ulid ulid.IULID,
	clock clock.IClock,
	transaction database.ITransaction,
	stateCountRepository repository.IStateCountRepository,
	subscribeStripeRepository repository.ISubscribeStripeRepository,
	userAccountLineRepository repository.IUserAccountLineRepository,
	userAggregateRepository repository.IUserAggregateRepository,
	userDetailRepository repository.IUserDetailRepository,
	userRepository repository.IUserRepository,
	llmClient client.ILLMClient,
	stripeDriver payment.IStripeDriver,
) IStripeSubscribeUseCase {
	return &StripeSubscribeUseCase{
		ulid:                      ulid,
		clock:                     clock,
		transaction:               transaction,
		stateCountRepository:      stateCountRepository,
		subscribeStripeRepository: subscribeStripeRepository,
		userAccountLineRepository: userAccountLineRepository,
		userAggregateRepository:   userAggregateRepository,
		userDetailRepository:      userDetailRepository,
		userRepository:            userRepository,
		llmClient:                 llmClient,
		stripeDriver:              stripeDriver,
	}
}

func (u *StripeSubscribeUseCase) Execute(req dto.Subscribe) (dto.SubscribeInfo, error) {
	// ユーザー取得または新規作成
	user, err := u.getUser(req.Account)
	if err != nil {
		return dto.SubscribeInfo{}, fmt.Errorf("failed to get user: %w", err)
	}

	// 現在のサブスク状態を確認
	subscribed, err := u.checkSubscribed(user.UserID)
	if err != nil {
		return dto.SubscribeInfo{}, fmt.Errorf("failed to check subscribed: %w", err)
	}
	if subscribed {
		return dto.SubscribeInfo{
			Status: entity.SubscribeStatusActive,
		}, nil
	}

	// サブスク開始のためのセッション作成
	redirectUrl, sessionID, err := u.stripeDriver.CreateSubscriptionSession(req)
	if err != nil {
		return dto.SubscribeInfo{}, fmt.Errorf("failed to create subscription session: %w", err)
	}

	// サブスク情報を作成
	if err := u.createSubscribeStripe(user.UserID, sessionID); err != nil {
		return dto.SubscribeInfo{}, fmt.Errorf("failed to create subscribe stripe: %w", err)
	}

	return dto.SubscribeInfo{
		RedirectUrl: redirectUrl,
	}, nil
}

func (u *StripeSubscribeUseCase) createSubscribeStripe(userID, sessionID string) error {
	subscribeStripe := entity.SubscribeStripe{
		SubscribeStripeID: u.ulid.GenerateID(),
		User: entity.User{
			UserID: userID,
		},
		SessionID: sessionID,
		Status:    entity.SubscribeStatusProcessing,
	}
	return u.subscribeStripeRepository.Create(subscribeStripe)
}

func (u *StripeSubscribeUseCase) getUser(account dto.Account) (entity.User, error) {
	user, err := u.userRepository.FindByAccount(account.PlatformType.String(), account.AccountID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// ユーザーが存在しない場合は新規作成
			if err := u.transaction.StartTransaction(func(tx any) error {
				var createErr error
				user, createErr = u.newUser(tx, account)
				if createErr != nil {
					return fmt.Errorf("failed to create user: %w", createErr)
				}
				return nil
			}); err != nil {
				return entity.User{}, fmt.Errorf("failed to start transaction for new user: %w", err)
			}
		} else {
			return entity.User{}, err
		}
	}
	return user, nil
}

func (u *StripeSubscribeUseCase) newUser(tx any, account dto.Account) (entity.User, error) {
	// ユーザー作成
	user, err := u.createUser(tx)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	// ユーザー詳細の作成
	if _, err = u.createUserDetail(tx, user.UserID); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user detail: %w", err)
	}
	// 状態カウントの作成
	if _, err = u.createStateCount(tx, user.UserID); err != nil {
		return entity.User{}, fmt.Errorf("failed to create state count: %w", err)
	}
	return user, nil
}

func (u *StripeSubscribeUseCase) createUser(tx any) (entity.User, error) {
	newUser := entity.User{
		UserID: u.ulid.GenerateID(),
	}
	if err := u.userRepository.CreateInTx(tx, newUser); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}

func (u *StripeSubscribeUseCase) createUserDetail(tx any, userID string) (entity.UserDetail, error) {
	userDetail := entity.UserDetail{
		UserDetailID: u.ulid.GenerateID(),
		User: entity.User{
			UserID: userID,
		},
	}
	if err := u.userDetailRepository.CreateInTx(tx, userDetail); err != nil {
		return entity.UserDetail{}, fmt.Errorf("failed to create user detail: %w", err)
	}
	return userDetail, nil
}

func (u *StripeSubscribeUseCase) createStateCount(tx any, userID string) (entity.StateCount, error) {
	stateCount := entity.StateCount{
		StateCountID: u.ulid.GenerateID(),
		User: entity.User{
			UserID: userID,
		},
		Count: 0,
	}
	if err := u.stateCountRepository.CreateInTx(tx, stateCount); err != nil {
		return entity.StateCount{}, fmt.Errorf("failed to create state count: %w", err)
	}
	return stateCount, nil
}

func (u *StripeSubscribeUseCase) checkSubscribed(userID string) (bool, error) {
	// サブスクステータスを最新に更新
	if err := u.updateSubscribeStatus(userID); err != nil {
		return false, fmt.Errorf("failed to update subscribe status: %w", err)
	}
	// サブスク情報を取得
	subscribeStripe, err := u.subscribeStripeRepository.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find subscribe stripe: %w", err)
	}
	return subscribeStripe.IsSubscribed(), nil
}

func (u *StripeSubscribeUseCase) updateSubscribeStatus(userID string) error {
	subscribeStripe, err := u.subscribeStripeRepository.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to find subscribe stripe: %w", err)
	}
	now := u.clock.Now().Unix()
	subscribeStatus, _, _, err := u.stripeDriver.GetSubscriptionStatus(subscribeStripe.SessionID, now)
	if err != nil {
		return fmt.Errorf("failed to get subscription status: %w", err)
	}
	subscribeStripe.Status = subscribeStatus
	if err := u.subscribeStripeRepository.Update(subscribeStripe); err != nil {
		return fmt.Errorf("failed to update subscribe stripe: %w", err)
	}
	return nil
}
