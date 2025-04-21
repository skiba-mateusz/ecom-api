package handler

import (
	"go.uber.org/zap"
	"net/http"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error, logger *zap.SugaredLogger) {
	logger.Errorw("internal server error", "path", r.URL.Path, "method", r.Method, "error", err.Error())
	_ = jsonErrorResponse(w, http.StatusInternalServerError, "internal server error")
}
