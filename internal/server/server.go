package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/romutchio/crypto-calculator/docs"

	"github.com/romutchio/crypto-calculator/api/routes"
	"github.com/romutchio/crypto-calculator/internal/usecase"
	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	useCase usecase.CalculateUseCase
}

func NewServer(calculateUc usecase.CalculateUseCase) *Server {
	app := fiber.New()
	app.Use(cors.New())

	return &Server{
		app:     app,
		useCase: calculateUc,
	}
}

func (s *Server) setupRoutes(ctx context.Context) {
	routes.HealthRouter(s.app)
	s.app.Get("/swagger/*", swagger.HandlerDefault)

	api := s.app.Group("/api")
	v1 := api.Group("/v1")
	routes.CalculatorRouter(ctx, v1, s.useCase)
}

func (s *Server) Run(ctx context.Context) error {
	// setup routes
	s.setupRoutes(ctx)
	// run http server

	chanErr := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting server on :8080")
		err := s.app.Listen(":8080")
		if err == nil {
			chanErr <- nil
			return
		}
		chanErr <- fmt.Errorf("s.app.Listen: %w", err)
	}()

	// waiting error or shutdown signal from <-ctx.Done()

	select {
	case err := <-chanErr:
		return err // when graceful got error then all runners got stop signal.

	case <-ctx.Done():
		log.Info().Msg("got shutdown signal")

		err := s.app.ShutdownWithTimeout(time.Minute)
		if err != nil {
			return fmt.Errorf("s.app.ShutdownWithTimeout: %w", err)
		}
		log.Info().Msg("http server graceful shutdown")
	}

	return nil
}
