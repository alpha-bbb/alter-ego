package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alpha-bbb/alter-ego/backend/application/dto"
	"github.com/alpha-bbb/alter-ego/backend/config"
	"github.com/alpha-bbb/alter-ego/backend/entity"
	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.uber.org/zap"
)

type StripeWebhookHandler struct {
	stripeSubscribeWebHookUseCase usecase.IStripeSubscribeWebhookUseCase
}

func NewStripeWebhookHandler(
	stripeSubscribeWebhookUseCase usecase.IStripeSubscribeWebhookUseCase,
) *StripeWebhookHandler {
	return &StripeWebhookHandler{
		stripeSubscribeWebHookUseCase: stripeSubscribeWebhookUseCase,
	}
}

func (h *StripeWebhookHandler) Execute(c echo.Context) error {
	logger, err := log.NewLogger()
	logger.Info("Stripe webhook handler")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create logger"})
	}

	const MaxBodyBytes = int64(65536)
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, MaxBodyBytes)

	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		logger.Error("Failed to read request body", zap.Error(err))
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "failed to read request body"})
	}

	sigHeader := c.Request().Header.Get("Stripe-Signature")
	if sigHeader == "" {
		logger.Error("Missing Stripe-Signature header")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing Stripe-Signature header"})
	}

	endpointSecret := config.StripeEndpointSecret()

	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		logger.Error("Failed to verify webhook signature", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "webhook signature verification failed"})
	}

	// イベントのタイプに応じた処理
	switch event.Type {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			logger.Error("Failed to parse event data", zap.Error(err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to parse event data"})
		}
		logger.Info("Checkout session completed", zap.String("session_id", sess.ID))
		sessionID := sess.ID
		email := sess.CustomerEmail
		return h.stripeSubscribeWebHookUseCase.Execute(
			dto.StripeSubscribeWebhook{
				SessionID: sessionID,
				Email:     email,
				Status:    entity.SubscribeStatusActive,
			},
		)
	default:
		logger.Info("Unhandled event type", zap.String("type", string(event.Type)))
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}
