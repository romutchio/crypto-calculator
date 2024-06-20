package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romutchio/crypto-calculator/internal/entity"
)

func (r *Repo) GetPairs(ctx context.Context) ([]*entity.Pair, error) {
	var pairs []*entity.Pair

	err := r.q(ctx).Find(&pairs).Error
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	return pairs, nil
}

func (r *Repo) CreatePair(ctx context.Context, pair *entity.Pair) error {
	err := r.q(ctx).Create(&pair).Error
	if err != nil {
		return errors.Wrap(err, "query error")
	}
	return nil
}

func (r *Repo) GetPair(ctx context.Context, from string, to string) (*entity.Pair, error) {
	var pair *entity.Pair

	err := r.q(ctx).Where("\"from\" = ? AND \"to\" = ?", from, to).Order("created_at desc").First(&pair).Error
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	return pair, nil
}
