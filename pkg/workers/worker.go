package workers

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/romutchio/crypto-calculator/internal/config"
	"github.com/romutchio/crypto-calculator/internal/entity"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetConfigurations(ctx context.Context) ([]*entity.Configuration, error)
	CreatePair(ctx context.Context, pair *entity.Pair) error
}

type DataProvider interface {
	FetchOne(from string, to string) (*entity.FetchOneResponse, error)
}

type PriceUpdater struct {
	cfg          *config.PriceUpdater
	repo         Repository
	dataProvider DataProvider
}

func NewPriceUpdater(cfg *config.PriceUpdater, repo Repository, dataProvider DataProvider) *PriceUpdater {
	return &PriceUpdater{
		cfg:          cfg,
		repo:         repo,
		dataProvider: dataProvider,
	}
}

type Request struct {
	From string
	To   string
}

func (pu *PriceUpdater) fetchData(requests <-chan Request, pairs chan<- *entity.Pair, wg *sync.WaitGroup) {
	for request := range requests {
		resp, err := pu.dataProvider.FetchOne(request.From, request.To)
		if err != nil {
			log.Error().Err(err).Msg("pu.dataProvider.FetchOne failed")
			continue
		}
		for to, amount := range resp.Result {
			pairs <- &entity.Pair{
				From:      resp.Base,
				To:        to,
				Amount:    amount,
				CreatedAt: time.Now().UTC(),
			}
		}
		wg.Done()
	}
}

func (pu *PriceUpdater) update(ctx context.Context) error {
	log.Info().Msg("start updating currencies")
	// Получаем список всех валют, которые мы должны поддерживать
	configurations, err := pu.repo.GetConfigurations(ctx)
	if err != nil {
		return errors.Wrap(err, "pu.repo.GetConfigurations failed")
	}
	if len(configurations) == 0 {
		log.Info().Msg("Price update skipped: no configurations")
		return nil
	}
	fiatSupportedCurrencies := make([]string, 0, len(configurations))
	cryptoSupportedCurrencies := make([]string, 0, len(configurations))
	for _, conf := range configurations {
		if conf.Type == entity.FIAT {
			fiatSupportedCurrencies = append(fiatSupportedCurrencies, conf.Code)
		}
		if conf.Type == entity.CRYPTO {
			cryptoSupportedCurrencies = append(cryptoSupportedCurrencies, conf.Code)
		}
	}

	if len(fiatSupportedCurrencies) == 0 {
		log.Info().Msg("Price update skipped: no fiat currencies supported in config")
		return nil
	}

	if len(cryptoSupportedCurrencies) == 0 || len(fiatSupportedCurrencies) == 0 {
		log.Info().Msg("Price update skipped: no crypto currencies supported in config")
		return nil
	}

	numWorkers := 4
	var wg sync.WaitGroup
	requests := make(chan Request, len(fiatSupportedCurrencies)*len(cryptoSupportedCurrencies))
	pairs := make(chan *entity.Pair, len(fiatSupportedCurrencies)*len(cryptoSupportedCurrencies))

	for w := 0; w < numWorkers; w++ {
		go pu.fetchData(requests, pairs, &wg)
	}
	wg.Add(len(fiatSupportedCurrencies) * len(cryptoSupportedCurrencies))
	for _, fromCurr := range fiatSupportedCurrencies {
		for _, toCurr := range cryptoSupportedCurrencies {
			requests <- Request{From: fromCurr, To: toCurr}
		}
	}
	close(requests)
	wg.Wait()

	// Сохраняем стоимость крипто пар
	for i := 0; i < len(fiatSupportedCurrencies)*len(cryptoSupportedCurrencies); i++ {
		err := pu.repo.CreatePair(ctx, <-pairs)
		if err != nil {
			return errors.Wrap(err, "pu.repo.CreatePair failed")
		}
	}
	log.Info().Msg("currencies updates successfully")
	return nil
}

func (pu *PriceUpdater) Run(ctx context.Context) error {
	ticker := time.NewTicker(pu.cfg.Period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("background work gracefully stopped")
			return nil

		case <-ticker.C:
			err := pu.update(ctx)
			if err != nil {
				log.Error().Err(err).Msg("background work has failed")
			}
		}
	}
}
