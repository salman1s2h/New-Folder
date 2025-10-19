package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true" env-default:"developm"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-required:"true" env-default:"8080"`
	Host string `yaml:"host" env:"HTTP_HOST" env-required:"true"  env-default:"localhost"`
}

func MustLoad() *Config {

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "./configs/config.yaml Path to config file")
		flag.Parse()
		configPath = *flags
	}
	if configPath == "" {
		log.Fatal("config path is required")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return &cfg
}
