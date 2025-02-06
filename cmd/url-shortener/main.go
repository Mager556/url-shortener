package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Mager556/url-shortener/internal/config"
	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/delete"
	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Mager556/url-shortener/internal/http-server/middleware/logger"
	"github.com/Mager556/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/Mager556/url-shortener/internal/lib/logger/sl"
	"github.com/Mager556/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/", delete.New(log, storage))
	})
	router.Get("/{alias}", redirect.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case EnvLocal:
		log = setupPrettySlog()
	case EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case EnvProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
