package main

import (
	"fmt"
	"github.com/idsulik/url-shortener/internal/config"
	"github.com/idsulik/url-shortener/internal/logger"
	"github.com/idsulik/url-shortener/internal/storage/sqlite"
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

	if err != nil {
		log.Error(fmt.Sprintf("Failed to initialize storage: %s", err.Error()))
		os.Exit(1)
	}

	log.Info("Starting url-shortener")
	log.Debug(fmt.Sprintf("Config: %+v", cfg))
	log.Debug(fmt.Sprintf("Storage: %+v", storage))
}
