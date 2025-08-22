package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	TelegramBotToken string
	TelegramGroupID  int64
	DatabasePath     string
	ProblemsFilePath string
	Timezone         string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramGroupID:  getEnvInt64("TELEGRAM_GROUP_ID", 0),
		DatabasePath:     getEnv("DATABASE_PATH", "leetcode_bot.db"),
		ProblemsFilePath: getEnv("PROBLEMS_FILE_PATH", "problem_deduplicated.yaml"),
		Timezone:         getEnv("TIMEZONE", "Asia/Ho_Chi_Minh"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
