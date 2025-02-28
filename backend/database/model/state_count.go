package model

import (
	"time"

	"github.com/alpha-bbb/alter-ego/backend/entity"
)

type StateCount struct {
	StateCountID string `gorm:"primaryKey;size:26"`
	UserID       string `gorm:"size:26"`
	User         User
	Count        int
	UpdatedAt    time.Time
}

func (s *StateCount) Entity() entity.StateCount {
	return entity.StateCount{
		StateCountID: s.StateCountID,
		User:         s.User.Entity(),
		Count:        s.Count,
		UpdatedAt:    s.UpdatedAt,
	}
}
