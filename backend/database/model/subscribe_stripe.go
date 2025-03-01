package model

import (
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type SubscribeStripe struct {
	SubscribeStripeID string `gorm:"primaryKey;size:26"`
	UserID            string `gorm:"size:26"`
	User              User
	SessionID         string `gorm:"size:255"`
	Status            string `gorm:"size:255"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

func (s *SubscribeStripe) Entity() entity.SubscribeStripe {
	status, _ := entity.ParseSubscribeStatus(s.Status)
	return entity.SubscribeStripe{
		SubscribeStripeID: s.SubscribeStripeID,
		User:              s.User.Entity(),
		SessionID:         s.SessionID,
		Status:            status,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}
