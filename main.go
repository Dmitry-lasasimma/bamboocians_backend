package main

import (
	"bamboocians/config"
	"bamboocians/database"
	"bamboocians/middleware"
	"bamboocians/routes"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	// Share JWT secret with auth middleware
	middleware.JWTSecret = cfg.JWTSecret

	database.Connect(cfg.DBPath)
	if cfg.SeedDemo {
		database.SeedDemo()
	}

	app := fiber.New(fiber.Config{
		AppName: "Bamboocians API v1.0",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	allowedOrigins := map[string]bool{}
	for _, o := range strings.Split(cfg.CORSOrigins, ",") {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			allowedOrigins[o] = true
		}
	}
	app.Use(cors.New(cors.Config{
		// Log rejected origins so a CORS_ORIGINS mismatch is visible in the server logs
		AllowOriginsFunc: func(origin string) bool {
			if allowedOrigins[origin] {
				return true
			}
			log.Printf("CORS blocked origin %q: add it to CORS_ORIGINS (currently %q)", origin, cfg.CORSOrigins)
			return false
		},
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	routes.Setup(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "app": "Bamboocians"})
	})

	log.Printf("Bamboocians API running on :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
