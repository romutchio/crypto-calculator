package entity

import (
	"time"
)

type ConfigurationType int16

const (
	FIAT ConfigurationType = iota
	CRYPTO
)

type Configuration struct {
	ID          uint              `gorm:"primarykey"`
	Name        string            `gorm:"column:name"`
	Code        string            `gorm:"column:code"`
	IsAvailable bool              `gorm:"column:is_available"`
	Type        ConfigurationType `gorm:"column:type"`
	CreatedAt   time.Time         `gorm:"column:created_at"`
}

func (*Configuration) TableName() string {
	return "configurations"
}

type Pair struct {
	ID        uint      `gorm:"primarykey"`
	From      string    `gorm:"column:from"`
	To        string    `gorm:"column:to"`
	Amount    float64   `gorm:"column:amount"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (*Pair) TableName() string {
	return "pairs"
}
