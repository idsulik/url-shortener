package main

import (
	"fmt"
	"github.com/idsulik/url-shortener/internal/config"
	"github.com/idsulik/url-shortener/internal/logger"
	"os"
)

func main() {
	env := os.Getenv("ENV")

	if env == "" {
		panic("ENV is not set")
	}

	cfg := config.New(env)
	log := logger.New(env)

	log.Info("Starting url-shortener")
	log.Debug(fmt.Sprintf("Config: %+v", cfg))
}
