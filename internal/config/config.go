package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	FastForex    FastForex
	DB           DB
	PriceUpdater PriceUpdater

	Help bool `default:"false"`
}

type FastForex struct {
	URL     string        `default:"https://api.fastforex.io/"`
	Token   string        `default:""`
	IsLocal bool          `default:"false" usage:"enable for unit tests"`
	Timeout time.Duration `default:"30s" usage:"timeout for requests"`
}

type DB struct {
	Host            string `default:"localhost" usage:"Host of the DB connect"`
	Port            int    `default:"5432" usage:"Port of the DB connect"`
	Name            string `usage:"Name of the DB"`
	User            string `usage:"Username of the connect"`
	Pass            string `usage:"Password of the connect" validate:"required"`
	Timeout         string `default:"60" usage:"Timeout of the DB connect"`
	SSL             string `usage:"Enables ssl connection to DB"`
	Timezone        string `env:"TZ" default:"Europe/Moscow" usage:"Timezone of the service"`
	ApplicationName string `default:"calculator" usage:"Application name of the DB connect"`

	MaxOpenConns int `default:"20" usage:"Max count of open connect in pool of DB"`
	MaxIdleConns int `default:"10" usage:"Max cout of idle connect in pool of DB"`
}

// DSN формирует строку подключения к базе
func (c *DB) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s  connect_timeout=%s TimeZone=%s  application_name=%s",
		c.Host, c.Port, c.User, c.Name, c.SSL, c.Pass, c.Timeout, c.Timezone, c.ApplicationName)
}

type PriceUpdater struct {
	Period time.Duration `default:"1m" usage:"price update period"`
}

func Load() (*Config, error) {
	conf := &Config{}
	loader := aconfig.LoaderFor(conf, aconfig.Config{})

	if err := loader.Load(); err != nil {
		return nil, fmt.Errorf("config load failed %w", err)
	}

	if conf.Help {
		loader.Flags().PrintDefaults()
		return nil, flag.ErrHelp
	}

	return conf, nil
}
