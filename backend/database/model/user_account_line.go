package model

import (
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
	"gorm.io/gorm"
)

type UserAccountLine struct {
	UserAccountLineID string `gorm:"primaryKey;size:26"`
	UserID            string `gorm:"size:26"`
	User              User
	LineID            string `gorm:"size:255"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

func (u *UserAccountLine) Entity() entity.UserAccountLine {
	return entity.UserAccountLine{
		UserAccountLineID: u.UserAccountLineID,
		User:              u.User.Entity(),
		LineID:            u.LineID,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}
