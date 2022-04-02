package toutcalc

import (
	"context"

	"gopkg.in/validator.v2"
)

type Calculator struct {
	easyTimeout        int
	mediumTimeout      int
	hardTimeout        int
	mediumCorrectCount int
	hardCorrectCount   int
}

type Config struct {
	EasyTimeout        int `validate:"nonzero"`
	MediumTimeout      int `validate:"nonzero"`
	HardTimeout        int `validate:"nonzero"`
	MediumCorrectCount int `validate:"nonzero"`
	HardCorrectCount   int `validate:"nonzero"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func StandardConfig() Config {
	return Config{
		EasyTimeout:        10,
		MediumTimeout:      5,
		HardTimeout:        3,
		MediumCorrectCount: 10,
		HardCorrectCount:   25,
	}
}

func New(cfg Config) (*Calculator, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	c := &Calculator{
		easyTimeout:        cfg.EasyTimeout,
		mediumTimeout:      cfg.MediumTimeout,
		hardTimeout:        cfg.HardTimeout,
		mediumCorrectCount: cfg.MediumCorrectCount,
		hardCorrectCount:   cfg.HardCorrectCount,
	}
	return c, nil
}

func (c *Calculator) CalcTimeout(ctx context.Context, countCorrect int) int {
	if countCorrect > c.hardCorrectCount {
		return c.hardTimeout
	}
	if countCorrect > c.mediumCorrectCount {
		return c.mediumTimeout
	}
	return c.easyTimeout
}
