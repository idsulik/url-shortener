package main

import (
	"fmt"
	"github.com/idsulik/url-shortener/internal/config"
	"os"
)

func main() {
	env := os.Getenv("ENV")

	if env == "" {
		panic("ENV is not set")
	}

	cfg := config.New(env)

	fmt.Println(cfg)
}
