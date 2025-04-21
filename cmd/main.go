package main

import (
	"github.com/skiba-mateusz/ecom-api/internal/infra/config"
	"github.com/skiba-mateusz/ecom-api/internal/infra/http"
	"github.com/skiba-mateusz/ecom-api/internal/infra/http/handler"
	"github.com/skiba-mateusz/ecom-api/internal/infra/persistance/postgres"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger := zap.Must(zap.NewProduction()).Sugar()

	db, err := postgres.New(
		cfg.Database.Addr,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.MaxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	handlers := &handler.Handlers{
		Health: handler.NewHealthHandler(cfg, logger),
	}

	server := http.NewServer(cfg, logger, handlers)
	mux := server.Mount()
	logger.Fatal(server.Run(mux))
}
