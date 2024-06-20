package usecase

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romutchio/crypto-calculator/internal/entity"
)

type Repository interface {
	GetPair(ctx context.Context, from string, to string) (*entity.Pair, error)
	GetConfigurationByCode(ctx context.Context, code string) (*entity.Configuration, error)
}

type CalculateUseCase struct {
	repo Repository
}

func NewCalculateUseCase(repo Repository) CalculateUseCase {
	return CalculateUseCase{
		repo: repo,
	}
}

type CalculateRequest struct {
	From   string
	To     string
	Amount float64
}

type CalculateResponse struct {
	Rate   float64
	Amount float64
}

func (uc *CalculateUseCase) Calculate(ctx context.Context, req *CalculateRequest) (*CalculateResponse, error) {
	fromConfiguration, err := uc.repo.GetConfigurationByCode(ctx, req.From)
	if err != nil {
		return nil, errors.Wrap(err, "uc.repo.GetConfiguration")
	}
	toConfiguration, err := uc.repo.GetConfigurationByCode(ctx, req.To)
	if err != nil {
		return nil, errors.Wrap(err, "uc.repo.GetConfiguration")
	}
	// Не нашли в конфигурации валюту
	if fromConfiguration == nil {
		return nil, errors.Errorf("unsupported currency 'from': %s", req.From)
	}
	if toConfiguration == nil {
		return nil, errors.Errorf("unsupported currency 'to': %s", req.To)
	}
	// Проверяем на включенность конфигурации
	if !fromConfiguration.IsAvailable {
		return nil, errors.Errorf("configuration is not available for currency: %s", req.From)
	}
	if !toConfiguration.IsAvailable {
		return nil, errors.Errorf("configuration is not available for currency: %s", req.To)
	}
	// Проверяем, что перевод C2F или F2C
	if fromConfiguration.Type == toConfiguration.Type {
		return nil, errors.New("only C2F and F2C calculation available")
	}
	rate := 1.0
	// В базе храним переходы fiat/crypto
	if fromConfiguration.Type == entity.CRYPTO && toConfiguration.Type == entity.FIAT {
		pair, err := uc.repo.GetPair(ctx, toConfiguration.Code, fromConfiguration.Code)
		if err != nil {
			return nil, errors.Wrap(err, "uc.repo.GetPair failed")
		}
		if pair == nil {
			return nil, errors.Errorf("pair was not found for %s/%s", toConfiguration.Code, fromConfiguration.Code)
		}
		if pair.Amount == 0 {
			return nil, errors.Errorf("currency rate can not be 0")
		}
		rate = 1 / pair.Amount
	}
	if fromConfiguration.Type == entity.FIAT && toConfiguration.Type == entity.CRYPTO {
		pair, err := uc.repo.GetPair(ctx, fromConfiguration.Code, toConfiguration.Code)
		if err != nil {
			return nil, errors.Wrap(err, "uc.repo.GetPair failed")
		}
		if pair == nil {
			return nil, errors.Errorf("pair was not found for %s/%s", fromConfiguration.Code, toConfiguration.Code)
		}
		if pair.Amount == 0 {
			return nil, errors.Errorf("currency rate can not be 0")
		}
		rate = pair.Amount
	}
	return &CalculateResponse{Amount: req.Amount * rate, Rate: rate}, nil
}
