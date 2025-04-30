package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/skiba-mateusz/ecom-api/internal/infra/config"
	"github.com/skiba-mateusz/ecom-api/internal/infra/http/handler"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	config   *config.Config
	logger   *zap.SugaredLogger
	handlers *handler.Handlers
}

func NewServer(config *config.Config, logger *zap.SugaredLogger, handlers *handler.Handlers) *Server {
	return &Server{
		config:   config,
		logger:   logger,
		handlers: handlers,
	}
}

func (s *Server) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", s.handlers.Health.CheckHealth)

		r.Route("/products", func(r chi.Router) {
			r.Get("/", s.handlers.Product.ListProducts)
			r.Post("/", s.handlers.Product.CreateProduct)

			r.Route("/{id}", func(r chi.Router) {
				r.Use(s.handlers.Product.ProductIdMiddleware)

				r.Get("/", s.handlers.Product.GetProduct)
				r.Put("/", s.handlers.Product.UpdateProduct)
				r.Delete("/", s.handlers.Product.DeleteProduct)
			})

		})
	})

	return r
}

func (s *Server) Run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         s.config.Http.Addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	s.logger.Infow("starting http server", "addr", s.config.Http.Addr, "env", s.config.Env)

	return srv.ListenAndServe()
}
