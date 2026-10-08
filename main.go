package main

import (
	"bamboocians/config"
	"bamboocians/database"
	"bamboocians/middleware"
	"bamboocians/routes"
	"log"

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

	app := fiber.New(fiber.Config{
		AppName: "Bamboocians API v1.0",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
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
