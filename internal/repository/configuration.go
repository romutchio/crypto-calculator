package repository

import (
	"context"

	"github.com/pkg/errors"

	"github.com/romutchio/crypto-calculator/internal/entity"
)

func (r *Repo) GetConfigurations(ctx context.Context) ([]*entity.Configuration, error) {
	var configurations []*entity.Configuration

	err := r.q(ctx).Find(&configurations).Error
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	return configurations, nil
}

func (r *Repo) GetConfigurationByCode(ctx context.Context, code string) (*entity.Configuration, error) {
	var configuration *entity.Configuration

	err := r.q(ctx).Where("code = ?", code).Find(&configuration).Error
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	return configuration, nil
}
