package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/romutchio/crypto-calculator/api/handlers"
)

func HealthRouter(app fiber.Router) fiber.Router {
	return app.Get("/", handlers.HealthCheck)
}
