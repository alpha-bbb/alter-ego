package entity

import "time"

type StateCount struct {
	StateCountID string
	User         User
	Count        int
	UpdatedAt    time.Time
}
