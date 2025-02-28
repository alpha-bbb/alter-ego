package repository

import (
	"errors"

	"github.com/alpha-bbb/alter-ego/backend/database/model"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type ISubscribeStripeRepository interface {
	Create(entity entity.SubscribeStripe) error
	FindByUserID(userID string) (entity.SubscribeStripe, error)
	Update(entity entity.SubscribeStripe) error
	UpdateStatusBySessionID(sessionID string, status entity.SubscribeStatus) error
}

type SubscribeStripeRepository struct {
	db *gorm.DB
}

func NewSubscribeStripeRepository(db *gorm.DB) ISubscribeStripeRepository {
	return &SubscribeStripeRepository{db: db}
}

func (r *SubscribeStripeRepository) Create(entity entity.SubscribeStripe) error {
	model := &model.SubscribeStripe{
		SubscribeStripeID: entity.SubscribeStripeID,
		UserID:            entity.User.UserID,
		SessionID:         entity.SessionID,
		Status:            entity.Status.String(),
	}
	return r.db.
		Create(model).
		Error
}
func (r *SubscribeStripeRepository) FindByUserID(userID string) (entity.SubscribeStripe, error) {
	var subscribeStripe model.SubscribeStripe
	err := r.db.Model(&model.SubscribeStripe{}).
		Preload("User").
		Where("user_id = ?", userID).
		First(&subscribeStripe).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.SubscribeStripe{}, ErrNotFound
		}
		return entity.SubscribeStripe{}, err
	}
	return subscribeStripe.Entity(), nil
}

func (r *SubscribeStripeRepository) Update(entity entity.SubscribeStripe) error {
	model := &model.SubscribeStripe{
		SubscribeStripeID: entity.SubscribeStripeID,
		UserID:            entity.User.UserID,
		SessionID:         entity.SessionID,
		Status:            entity.Status.String(),
	}
	return r.db.
		Save(model).
		Error
}

func (r *SubscribeStripeRepository) UpdateStatusBySessionID(sessionID string, status entity.SubscribeStatus) error {
	return r.db.
		Model(&model.SubscribeStripe{}).
		Where("session_id = ?", sessionID).
		Update("status", status.String()).
		Error
}
