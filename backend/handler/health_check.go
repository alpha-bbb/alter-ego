package handler

import (
	"net/http"

	"github.com/alpha-bbb/alter-ego/backend/infrastructure/log"
	"github.com/labstack/echo/v4"
)

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Execute(c echo.Context) error {
	logger, err := log.NewLogger()
	if err != nil {
		return err
	}

	logger.Info("Health check requested")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
