package handler

import (
	"go.uber.org/zap"
	"net/http"
)

func internalServerError(w http.ResponseWriter, r *http.Request, err error, logger *zap.SugaredLogger) {
	logger.Errorw("internal server error", "path", r.URL.Path, "method", r.Method, "error", err.Error())
	_ = jsonErrorResponse(w, http.StatusInternalServerError, "internal server error")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error, logger *zap.SugaredLogger) {
	logger.Warnw("bad request response", "path", r.URL.Path, "method", r.Method, "error", err.Error())
	_ = jsonErrorResponse(w, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request, err error, logger *zap.SugaredLogger) {
	logger.Warnw("not found response", "path", r.URL.Path, "method", r.Method, "error", err.Error())
	_ = jsonErrorResponse(w, http.StatusInternalServerError, "not found")
}
