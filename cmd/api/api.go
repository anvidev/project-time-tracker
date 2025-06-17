package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anvidev/project-time-tracker/internal/database"
	"github.com/anvidev/project-time-tracker/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *api) handler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", api.authRegister)
			r.Post("/login", api.authLogin)
		})

		r.Route("/time_entries", func(r chi.Router) {
			r.Use(api.bearerAuthorization)
			r.Get("/categories", api.entriesCategories)
		})
	})

	return r
}

type api struct {
	config config
	logger *slog.Logger
	store  *store.Store
}

func (api *api) Run() error {
	mux := api.handler()

	srv := &http.Server{
		Addr:         api.config.server.addr,
		ReadTimeout:  api.config.server.readTimeout,
		WriteTimeout: api.config.server.writeTimeout,
		IdleTimeout:  api.config.server.idleTimeout,
		Handler:      mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		api.logger.Info("server starting", "addr", api.config.server.addr, "env", api.config.server.env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			api.logger.Error("server failed to start", "error", err)
		}
	}()

	<-quit
	api.logger.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

func NewApiContext(ctx context.Context) (*api, error) {
	logger := slog.Default()
	config := loadConfig()

	db, err := database.NewContext(ctx, config.database.url, config.database.token)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		return nil, err
	}

	store := store.NewStore(db)

	api := &api{
		logger: logger,
		config: config,
		store:  store,
	}

	return api, nil
}
