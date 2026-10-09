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
	// Insert demo accounts and sample data on startup
	SeedDemo bool
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		JWTSecret:   getEnv("JWT_SECRET", "bamboocians-secret"),
		Port:        getEnv("PORT", "8081"),
		DBPath:      getEnv("DB_PATH", "./bamboocians.db"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
		SeedDemo:    getEnv("SEED_DEMO", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
