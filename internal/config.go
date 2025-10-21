package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string `yaml:"env" env:"ENV" env-default:"development"`
	HTTP struct {
		Port            string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
		Host            string `yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
		ShutdownTimeout int    `yaml:"ShutdownTimeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
	} `yaml:"http"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flag.StringVar(&configPath, "config", "./configs/config.yaml", "Path to config file")
		flag.Parse()
	}

	if configPath == "" {
		log.Fatal("config path is required")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}
