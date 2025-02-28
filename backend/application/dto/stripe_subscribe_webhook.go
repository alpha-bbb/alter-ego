package dto

import "github.com/alpha-bbb/alter-ego/backend/entity"

type StripeSubscribeWebhook struct {
	SessionID string
	Email     string
	Status    entity.SubscribeStatus
}
