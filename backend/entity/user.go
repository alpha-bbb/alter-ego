package entity

import "time"

type User struct {
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
