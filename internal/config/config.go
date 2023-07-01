package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type HttpServer struct {
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Config struct {
	Env         string
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func New(env string) *Config {
	configPath := "config/" + env + ".yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file '%s' does not exist", configPath))
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		panic(err)
	}

	config.Env = env

	return &config
}
