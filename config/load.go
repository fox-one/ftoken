package config

import (
	"github.com/fox-one/pkg/config"
	"github.com/shopspring/decimal"
)

func Load(cfgFile string, cfg *Config) error {
	config.AutomaticLoadEnv("FTOKEN")
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}

	defaultGas(cfg)

	return nil
}

func defaultGas(cfg *Config) {
	if cfg.Gas.StrictMultiplier.IsZero() {
		cfg.Gas.StrictMultiplier = decimal.New(4, 0)
	}

	if cfg.Gas.Multiplier.LessThan(cfg.Gas.StrictMultiplier) {
		cfg.Gas.Multiplier = cfg.Gas.StrictMultiplier.Add(decimal.New(1, 0))
	}
}
