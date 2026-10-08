package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	Port      string
	DBPath    string
	// Comma-separated list of allowed CORS origins
	CORSOrigins string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		JWTSecret:   getEnv("JWT_SECRET", "bamboocians-secret"),
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "./bamboocians.db"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
