package entity

import (
	"time"
)

type UserDetail struct {
	UserDetailID string
	User         User
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
