package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	DBPath           string
	LogLevel         string
}

func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	cfg := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		DBPath:           os.Getenv("DB_PATH"),
		LogLevel:         os.Getenv("LOG_LEVEL"),
	}

	// Validate required fields
	if cfg.TelegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN tidak ditemukan di .env")
	}

	// Set defaults
	if cfg.DBPath == "" {
		cfg.DBPath = "./data/schedules.json"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "INFO"
	}

	return cfg, nil
}
