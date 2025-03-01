package entity

import (
	"time"
)

type UserAccountLine struct {
	UserAccountLineID string
	User              User
	LineID            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
