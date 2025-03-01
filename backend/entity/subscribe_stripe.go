package entity

import (
	"fmt"
	"time"
)

type SubscribeStatus int

const (
	SubscribeStatusUnspecified   SubscribeStatus = iota // 未定義
	SubscribeStatusNotSubscribed                        // サブスクしていない状態
	SubscribeStatusProcessing                           // 決済処理中
	SubscribeStatusActive                               // サブスク中
)

func (s SubscribeStatus) String() string {
	switch s {
	case SubscribeStatusNotSubscribed:
		return "not_subscribed"
	case SubscribeStatusProcessing:
		return "processing"
	case SubscribeStatusActive:
		return "active"
	default:
		return "unspecified"
	}
}

func ParseSubscribeStatus(status string) (SubscribeStatus, error) {
	switch status {
	case "unspecified":
		return SubscribeStatusUnspecified, nil
	case "not_subscribed":
		return SubscribeStatusNotSubscribed, nil
	case "processing":
		return SubscribeStatusProcessing, nil
	case "active":
		return SubscribeStatusActive, nil
	default:
		return SubscribeStatusUnspecified, fmt.Errorf("invalid subscribe status: %s", status)
	}
}

type SubscribeStripe struct {
	SubscribeStripeID string
	User              User
	SessionID         string
	Status            SubscribeStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (e *SubscribeStripe) IsSubscribed() bool {
	return e.Status == SubscribeStatusActive
}
