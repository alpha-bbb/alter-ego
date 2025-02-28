package dto

import "github.com/alpha-bbb/alter-ego/backend/entity"

type SubscribeInfo struct {
	Status      entity.SubscribeStatus
	Message     string
	RedirectUrl string
	ExpiresAt   int64
}
