package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anvidev/apiduck"
	"github.com/anvidev/project-time-tracker/internal/database"
	"github.com/anvidev/project-time-tracker/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

func (api *api) handler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/docs", api.docs.Serve)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", api.authRegister)
			r.Post("/login", api.authLogin)
		})

		r.Route("/me", func(r chi.Router) {
			r.Use(api.bearerAuthorization)
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", api.entriesCategories)
				r.Get("/all", api.entriesCategoriesTree)
				r.Put("/{id}/follow", api.entriesFollowCategory)
				r.Put("/{id}/unfollow", api.entriesUnfollowCategory)
			})

			r.Route("/time_entries", func(r chi.Router) {
				r.Post("/", api.entriesRegisterTime)
				r.Put("/{id}", api.entriesUpdateTime)
				r.Delete("/{id}", api.entriesDelete)
				r.Get("/day/{date}", api.entriesSummaryDay)           // date: YYYY-MM-DD
				r.Get("/month/{year-month}", api.entriesSummaryMonth) // month: YYYY-MM
			})
		})

	})

	return r
}

type api struct {
	config config
	logger *slog.Logger
	store  *store.Store
	docs   *apiduck.Documentation
	cron   gocron.Scheduler
}

func (api *api) Run() error {
	mux := api.handler()

	if err := api.initCronScheduler(); err != nil {
		api.logger.Error("failed to init cron scheduler")
	}

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

	api.createCronJobs()
	go api.cron.Start()

	<-quit
	api.logger.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = api.cron.Shutdown()

	return srv.Shutdown(shutdownCtx)
}

func NewApiContext(ctx context.Context) (*api, error) {
	logger := slog.Default()
	config := loadConfig()
	docs := initDocumentation(config)

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
		docs:   docs,
	}

	return api, nil
}
