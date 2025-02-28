package model

import (
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type User struct {
	UserID    string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Detail          *UserDetail      `gorm:"foreignKey:UserID"`
	AccountLine     *UserAccountLine `gorm:"foreignKey:UserID"`
	SubscribeStripe *SubscribeStripe `gorm:"foreignKey:UserID"`
	StateCount      *StateCount      `gorm:"foreignKey:UserID"`
}

func (u *User) Entity() entity.User {
	return entity.User{
		UserID:    u.UserID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
