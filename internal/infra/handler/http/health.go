package http

import (
	"github.com/skiba-mateusz/ecom-api/internal/infra/config"
	"go.uber.org/zap"
	"net/http"
)

type HealthHandler struct {
	config *config.Config
	logger *zap.SugaredLogger
}

func NewHealthHandler(config *config.Config, logger *zap.SugaredLogger) *HealthHandler {
	return &HealthHandler{
		config: config,
		logger: logger,
	}
}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    h.config.Env,
	}

	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		internalServerError(w, r, err, h.logger)
	}
}
