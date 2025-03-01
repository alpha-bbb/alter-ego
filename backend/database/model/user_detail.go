package model

import (
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type UserDetail struct {
	UserDetailID string `gorm:"primaryKey;size:26"`
	UserID       string `gorm:"size:26"`
	User         User
	Email        string `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (u *UserDetail) Entity() entity.UserDetail {
	return entity.UserDetail{
		UserDetailID: u.UserDetailID,
		User:         u.User.Entity(),
		Email:        u.Email,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
