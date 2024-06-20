package routes

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/romutchio/crypto-calculator/api/handlers"
	"github.com/romutchio/crypto-calculator/internal/usecase"
)

func CalculatorRouter(ctx context.Context, app fiber.Router, uc usecase.CalculateUseCase) fiber.Router {
	return app.Get("/calculate", handlers.CalculateHandler(ctx, uc))
}
