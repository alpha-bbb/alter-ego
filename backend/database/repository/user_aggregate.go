package repository

import (
	"errors"

	"github.com/alpha-bbb/alter-ego/backend/database/model"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type IUserAggregateRepository interface {
	FindByUserID(userID string) (entity.UserAggregate, error)
}

type userAggregateRepository struct {
	db *gorm.DB
}

func NewUserAggregateRepository(db *gorm.DB) IUserAggregateRepository {
	return &userAggregateRepository{db: db}
}

func (r *userAggregateRepository) FindByUserID(userID string) (entity.UserAggregate, error) {
	var userModel model.User
	err := r.db.
		Where("user_id = ?", userID).
		Preload("Detail").
		Preload("AccountLine").
		Preload("SubscribeStripe").
		Preload("StateCount").
		First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.UserAggregate{}, ErrNotFound
		}
		return entity.UserAggregate{}, err
	}

	agg := entity.UserAggregate{
		User: userModel.Entity(),
	}

	if userModel.Detail != nil {
		d := userModel.Detail.Entity()
		agg.Detail = &d
	}
	if userModel.AccountLine != nil {
		ll := userModel.AccountLine.Entity()
		agg.AccountLine = &ll
	}
	if userModel.SubscribeStripe != nil {
		ss := userModel.SubscribeStripe.Entity()
		agg.StripeSubscription = &ss
	}
	if userModel.StateCount != nil {
		sc := userModel.StateCount.Entity()
		agg.StateCount = &sc
	}

	return agg, nil
}
