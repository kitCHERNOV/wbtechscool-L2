package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config структура
type Config struct {
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8081"`
}

func MustLoad() *Config {
	const op = "config.config.MustLoad"
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Loc: %s, Err: %v", op, err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
