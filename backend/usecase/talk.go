package usecase

import (
	"errors"
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/adapter/clock"
	"github.com/alpha-bbb/alter-ego/backend/adapter/llm/client"
	converterLlm "github.com/alpha-bbb/alter-ego/backend/adapter/llm/converter"
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/adapter/ulid"
	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/database"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/entity"
)

var ErrCountLimitExceeded = errors.New("count limit exceeded")

type ITalkUseCase interface {
	Execute(entity dto.Talk) (dto.Message, error)
}

type TalkUseCase struct {
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

func NewTalkUseCase(
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
) ITalkUseCase {
	return &TalkUseCase{
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

func (u *TalkUseCase) Execute(req dto.Talk) (dto.Message, error) {
	user, err := u.getUser(req.Account)
	if err != nil {
		return dto.Message{}, fmt.Errorf("failed to get user: %w", err)
	}
	subscribed, err := u.checkSubscribed(user.UserID)
	if err != nil {
		return dto.Message{}, fmt.Errorf("failed to check subscribed: %w", err)
	}
	if !subscribed {
		stateCount, err := u.stateCountRepository.FindByUserID(user.UserID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return dto.Message{}, fmt.Errorf("failed to find state count: %w", err)
		}
		countLimit := 3
		if stateCount.Count >= countLimit {
			return dto.Message{}, ErrCountLimitExceeded
		}
	}

	talkRequest := converterLlm.ToGRPCTalkRequest(req.Histories)
	var res dto.Message
	if err := u.transaction.StartTransaction(func(tx any) error {
		llmResponse, err := u.llmClient.Talk(talkRequest)
		if err != nil {
			return fmt.Errorf("failed to call LLM service: %w", err)
		}
		res = converterLlm.FromGRPCTalkResponse(llmResponse)
		return u.stateCountRepository.IncrementCountInTx(tx, user.UserID)
	}); err != nil {
		return dto.Message{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	return res, nil
}

func (u *TalkUseCase) getUser(account dto.Account) (entity.User, error) {
	user, err := u.userRepository.FindByAccount(account.PlatformType.String(), account.AccountID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return u.createUserWithTransaction(account)
		}
		return entity.User{}, err
	}
	return user, nil
}

func (u *TalkUseCase) createUserWithTransaction(account dto.Account) (entity.User, error) {
	var newUser entity.User
	if err := u.transaction.StartTransaction(func(tx any) error {
		var err error
		newUser, err = u.newUser(tx, account)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	}); err != nil {
		return entity.User{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	return newUser, nil
}

func (u *TalkUseCase) newUser(tx any, account dto.Account) (entity.User, error) {
	user, err := u.createUser(tx)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	// ユーザー詳細の作成
	if _, err = u.createUserDetail(tx, user.UserID); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user detail: %w", err)
	}
	// user_account_lines の作成を追加
	if err = u.createUserAccountLine(tx, account, user.UserID); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user account line: %w", err)
	}
	// 状態カウントの作成
	if _, err = u.createStateCount(tx, user.UserID); err != nil {
		return entity.User{}, fmt.Errorf("failed to create state count: %w", err)
	}
	return user, nil
}

func (u *TalkUseCase) createUser(tx any) (entity.User, error) {
	newUser := entity.User{
		UserID: u.ulid.GenerateID(),
	}
	if err := u.userRepository.CreateInTx(tx, newUser); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return newUser, nil
}

func (u *TalkUseCase) createUserDetail(tx any, userID string) (entity.UserDetail, error) {
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

func (u *TalkUseCase) createUserAccountLine(tx any, account dto.Account, userID string) error {
	// ユーザーアカウント連携情報を作成する
	userAccountLine := entity.UserAccountLine{
		UserAccountLineID: u.ulid.GenerateID(),
		User: entity.User{
			UserID: userID,
		},
		// 今回はLINEのみ対応している前提
		LineID: account.AccountID,
		// 必要に応じて他のフィールドも設定する
	}
	if err := u.userAccountLineRepository.CreateInTx(tx, userAccountLine); err != nil {
		return fmt.Errorf("failed to create user account line: %w", err)
	}
	return nil
}

func (u *TalkUseCase) createStateCount(tx any, userID string) (entity.StateCount, error) {
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

func (u *TalkUseCase) checkSubscribed(userID string) (bool, error) {
	if err := u.updateSubscribeStatus(userID); err != nil {
		return false, fmt.Errorf("failed to update subscribe status: %w", err)
	}
	subscribeStripe, err := u.subscribeStripeRepository.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to find subscribe stripe: %w", err)
	}
	return subscribeStripe.IsSubscribed(), nil
}

func (u *TalkUseCase) updateSubscribeStatus(userID string) error {
	subscribeStripe, err := u.subscribeStripeRepository.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to find subscribe stripe: %w", err)
	}
	now := u.clock.Now().Unix()
	subscribeStatus, _, err := u.stripeDriver.GetSubscriptionStatus(subscribeStripe.SessionID, now)
	if err != nil {
		return fmt.Errorf("failed to get subscription status: %w", err)
	}
	subscribeStripe.Status = subscribeStatus
	if err := u.subscribeStripeRepository.Update(subscribeStripe); err != nil {
		return fmt.Errorf("failed to update subscribe stripe: %w", err)
	}
	return nil
}
