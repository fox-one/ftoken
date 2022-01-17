package cmd

import (
	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/ftoken/quorum"
)

func provideEthFactory() core.Factory {
	return quorum.New(cfg.Eth.Endpoint, cfg.Eth.PrivateKey, cfg.Eth.FactoryContract)
}
