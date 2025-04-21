package main

import (
	"github.com/skiba-mateusz/ecom-api/internal/infra/config"
	"github.com/skiba-mateusz/ecom-api/internal/infra/http"
	"github.com/skiba-mateusz/ecom-api/internal/infra/http/handler"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger := zap.Must(zap.NewProduction()).Sugar()

	handlers := &handler.Handlers{
		Health: handler.NewHealthHandler(cfg, logger),
	}

	server := http.NewServer(cfg, logger, handlers)
	mux := server.Mount()
	logger.Fatal(server.Run(mux))
}
