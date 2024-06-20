package main

import (
	"flag"

	"github.com/pkg/errors"
	_ "github.com/romutchio/crypto-calculator/docs"
	"github.com/romutchio/crypto-calculator/internal/config"
	"github.com/romutchio/crypto-calculator/internal/repository"
	"github.com/romutchio/crypto-calculator/internal/server"
	"github.com/romutchio/crypto-calculator/internal/usecase"
	"github.com/romutchio/crypto-calculator/pkg/clients/fastforex"
	"github.com/romutchio/crypto-calculator/pkg/workers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/so_literate/graceful"
)

// @title			Calculator API
// @version		1.0
// @description	Calculator API Description
// @host			localhost:8082
// @BasePath		/
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatal().Err(err).Msg("config loading error")
	}

	fastForexClient, err := fastforex.New(&cfg.FastForex, nil)
	if err != nil {
		log.Error().Err(err).Msg("fastforex.New failed")
		return
	}

	repo, err := repository.New(&cfg.DB)
	if err != nil {
		log.Error().Err(err).Msg("repository.New failed")
		return
	}

	calculatorUc := usecase.NewCalculateUseCase(repo)
	httpServer := server.NewServer(calculatorUc)

	priceUpdater := workers.NewPriceUpdater(&cfg.PriceUpdater, repo, fastForexClient)

	err = graceful.Run(priceUpdater, httpServer)
	if err != nil {
		log.Error().Err(err).Msg("Shutdown with error")
		return
	}
}
