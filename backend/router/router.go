package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/alpha-bbb/alter-ego/backend/adapter/database/postgresql"
	"github.com/alpha-bbb/alter-ego/backend/adapter/payment"
	"github.com/alpha-bbb/alter-ego/backend/config"
	"github.com/alpha-bbb/alter-ego/backend/database/repository"
	"github.com/alpha-bbb/alter-ego/backend/handler"
	"github.com/alpha-bbb/alter-ego/backend/usecase"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := postgresql.NewPostgreSQL(config.DSN())
	if err != nil {
		panic(err)
	}

	stripeDriver := payment.NewStripeDriver(
		config.StripeApiKey(),
		config.StripeEndpointSecret(),
		config.FrontendURL(),
		config.PlaceID(),
	)

	subscribeStripeRepository := repository.NewSubscribeStripeRepository(db)

	stripeSubscribeWebhookUseCase := usecase.NewStripeSubscribeWebhookUseCase(
		subscribeStripeRepository,
		stripeDriver,
	)

	healthCheckHandler := handler.NewHealthCheckHandler()
	stripeWebhookHandler := handler.NewStripeWebhookHandler(
		stripeSubscribeWebhookUseCase,
	)

	e.GET("/health", healthCheckHandler.Execute)

	e.POST("/stripe/webhook", stripeWebhookHandler.Execute)

	return e
}
