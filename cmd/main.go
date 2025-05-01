package main

import (
	"github.com/skiba-mateusz/ecom-api/internal/app/service"
	"github.com/skiba-mateusz/ecom-api/internal/infra/config"
	"github.com/skiba-mateusz/ecom-api/internal/infra/handler/http"
	"github.com/skiba-mateusz/ecom-api/internal/infra/persistence/postgres"
	"github.com/skiba-mateusz/ecom-api/internal/infra/persistence/postgres/repository"
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

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productServ := service.NewProductService(productRepo, categoryRepo)

	handlers := &http.Handlers{
		Health:  http.NewHealthHandler(cfg, logger),
		Product: http.NewProductHandler(cfg, logger, productServ),
	}

	server := http.NewServer(cfg, logger, handlers)
	mux := server.Mount()
	logger.Fatal(server.Run(mux))
}
