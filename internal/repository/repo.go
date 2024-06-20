package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/romutchio/crypto-calculator/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repo struct {
	cfg *config.DB
	db  *gorm.DB
}

func New(cfg *config.DB) (*Repo, error) {
	d, err := gorm.Open(&postgres.Dialector{
		Config: &postgres.Config{
			DSN:                  cfg.DSN(),
			PreferSimpleProtocol: true,
		},
	}, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable:       false,
			NameReplacer:        nil,
			NoLowerCase:         false,
			IdentifierMaxLength: 0,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open failed")
	}

	db, err := d.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get *sql.DB")
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}

	return &Repo{cfg: cfg, db: d}, nil
}

func (r *Repo) q(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}
