package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	Port      string
	DBPath    string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		JWTSecret: getEnv("JWT_SECRET", "bamboocians-secret"),
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "./bamboocians.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
