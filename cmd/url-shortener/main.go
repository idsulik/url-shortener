package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idsulik/url-shortener/internal/alias"
	"github.com/idsulik/url-shortener/internal/config"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/save"
	loggermiddleware "github.com/idsulik/url-shortener/internal/http-server/middleware/logger-middleware"
	routerpackage "github.com/idsulik/url-shortener/internal/http-server/router"
	"github.com/idsulik/url-shortener/internal/logger"
	"github.com/idsulik/url-shortener/internal/storage/sqlite"
	"net/http"
	"os"
)

func main() {
	env := os.Getenv("ENV")

	if env == "" {
		panic("ENV is not set")
	}

	cfg := config.New(env)
	log := logger.New(env)
	storage, err := sqlite.New(cfg.StoragePath)
	aliasGenerator := alias.NewAliasGenerator()
	router := routerpackage.New()
	router.Use(loggermiddleware.New(log))
	router.Use(middleware.URLFormat)
	router.Route("/api", func(r chi.Router) {
		router.Route("/url", func(r chi.Router) {
			router.Use(
				middleware.BasicAuth(
					"url-shortener",
					map[string]string{cfg.HttpServer.User: cfg.HttpServer.Password},
				),
			)
			router.Post("/shorten", save.New(log, aliasGenerator, storage))
		})
	})
	router.Get("/{alias}", redirect.New(log, storage))

	if err != nil {
		log.Error(fmt.Sprintf("Failed to initialize storage: %s", err.Error()))
		os.Exit(1)
	}

	log.Info("Starting url-shortener")
	log.Debug(fmt.Sprintf("Config: %+v", cfg))
	log.Debug(fmt.Sprintf("Storage: %+v", storage))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HttpServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(fmt.Sprintf("Failed to start http server: %s", err.Error()))
		os.Exit(1)
	}
}
