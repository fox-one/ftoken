package config

import (
	"math/big"

	"github.com/fox-one/pkg/config"
)

func Load(cfgFile string, cfg *Config) error {
	config.AutomaticLoadEnv("FTOKEN")
	if err := config.LoadYaml(cfgFile, cfg); err != nil {
		return err
	}

	defaultEth(cfg)

	return nil
}

func defaultEth(cfg *Config) {
	if cfg.Eth.MaxGasPrice == nil {
		cfg.Eth.MaxGasPrice = big.NewInt(10000000000)
	}
}
