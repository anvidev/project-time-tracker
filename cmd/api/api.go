package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anvidev/apiduck"
	"github.com/anvidev/goenv"
	"github.com/anvidev/project-time-tracker/internal/database"
	"github.com/anvidev/project-time-tracker/internal/mailer"
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
				r.Post("/", api.entriesCreateCategory)
				r.Put("/{id}", api.entriesUpdateCategory)
				r.Put("/{id}/toggle", api.entriesToggleCategory)
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

			r.Route("/hours", func(r chi.Router) {
				r.Get("/", api.hoursAll)
			})
		})

	})

	return r
}

type api struct {
	config Config
	logger *slog.Logger
	store  *store.Store
	docs   *apiduck.Documentation
	mails  mailer.Mailer

	cronInitialized bool
	cron            gocron.Scheduler
}

func (api *api) Run() error {
	mux := api.handler()

	srv := &http.Server{
		Addr:         api.config.Server.Addr,
		ReadTimeout:  api.config.Server.ReadTimeout,
		WriteTimeout: api.config.Server.WriteTimeout,
		IdleTimeout:  api.config.Server.IdleTimeout,
		Handler:      mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		api.logger.Info("server starting", "addr", api.config.Server.Addr, "env", api.config.Server.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			api.logger.Error("server failed to start", "error", err)
		}
	}()

	if api.cronInitialized {
		if err := api.createCronJobs(); err != nil {
			api.logger.Warn("creating cron jobs failed", "error", err)
		} else {
			go api.cron.Start()
		}
	}

	<-quit
	api.logger.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = api.cron.Shutdown()

	return srv.Shutdown(shutdownCtx)
}

func NewApiContext(ctx context.Context) (*api, error) {
	logger := slog.Default()

	var config Config
	if err := goenv.Struct(&config); err != nil {
		logger.Error("error", "err", err.Error())
		os.Exit(1)
	}
	fmt.Printf("config is: %+v\n", config)

	docs := initDocumentation(config)
	mails := mailer.NewResendMailer(config.Resend.ApiKey, config.Resend.From)

	var cronInitialized bool
	cron, err := initCronScheduler()
	if err != nil {
		logger.Warn("cron scheduler initialization failed", "error", err)
	} else {
		cronInitialized = true
	}

	db, err := database.NewContext(ctx, config.Database.URL, config.Database.Token)
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
		mails:  mails,

		cronInitialized: cronInitialized,
		cron:            cron,
	}

	return api, nil
}
