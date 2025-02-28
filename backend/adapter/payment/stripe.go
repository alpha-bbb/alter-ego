package payment

import (
	"fmt"

	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/subscription"
)

type IStripeDriver interface {
	CreateSubscriptionSession(dto dto.Subscribe) (string, string, error)
	GetSubscriptionStatus(sessionID string, now int64) (entity.SubscribeStatus, int64, error)
	CancelSubscription(sessionID string) error
}

type StripeDriver struct {
	StripeEndpointSecret string
	FrontendURL          string
	PriceID              string
}

func NewStripeDriver(stripeApiKey, stripeEndpointSecret, frontendURL, priceID string) IStripeDriver {
	stripe.Key = stripeApiKey
	return &StripeDriver{
		StripeEndpointSecret: stripeEndpointSecret,
		FrontendURL:          frontendURL,
		PriceID:              priceID,
	}
}

func (s *StripeDriver) CreateSubscriptionSession(dto dto.Subscribe) (string, string, error) {
	logger, _ := log.NewLogger()
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("subscription"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(s.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(s.FrontendURL + "/success?session_id={CHECKOUT_SESSION_ID}"),
		Metadata: map[string]string{
			"platform_type": string(dto.Account.PlatformType),
			"account_id":    dto.Account.AccountID,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		logger.Info(fmt.Sprintf("failed to create session: %v", err))
		return "", "", err
	}
	return sess.URL, sess.ID, nil
}

func (s *StripeDriver) GetSubscriptionStatus(sessionID string, now int64) (entity.SubscribeStatus, int64, error) {
	logger, _ := log.NewLogger()
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("session.Get error: %v", err))
		return entity.SubscribeStatusUnspecified, 0, fmt.Errorf("failed to retrieve session: %w", err)
	}

	if sess.Subscription == nil {
		logger.Info(fmt.Sprintf("no subscription found in session: %+v", sess))
		return entity.SubscribeStatusUnspecified, 0, nil
	}
	subscriptionID := sess.Subscription.ID
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return entity.SubscribeStatusUnspecified, 0, fmt.Errorf("failed to retrieve subscription details: %w", err)
	}
	currentPeriodEnd := sub.CurrentPeriodEnd

	logger.Info(fmt.Sprintf("currentPeriodEnd: %d, now: %d", currentPeriodEnd, now))

	var status entity.SubscribeStatus
	if currentPeriodEnd < now {
		status = entity.SubscribeStatusNotSubscribed
	} else {
		status = entity.SubscribeStatusActive
	}
	return status, currentPeriodEnd, nil
}

func (s *StripeDriver) CancelSubscription(sessionID string) error {
	logger, _ := log.NewLogger()
	sess, err := session.Get(sessionID, nil)
	if err != nil {
		logger.Info(fmt.Sprintf("session.Get error: %v", err))
		return fmt.Errorf("failed to retrieve session: %w", err)
	}
	if sess.Subscription == nil {
		return fmt.Errorf("no subscription found for session %s", sessionID)
	}

	subscriptionID := sess.Subscription.ID
	cancelParams := &stripe.SubscriptionCancelParams{}
	cancelledSub, err := subscription.Cancel(subscriptionID, cancelParams)
	if err != nil {
		logger.Info(fmt.Sprintf("subscription.Cancel error: %v", err))
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}
	logger.Info(fmt.Sprintf("cancelled subscription: %+v", cancelledSub))
	return nil
}
