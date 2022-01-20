package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/shopspring/decimal"
)

type (
	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret"`
		Pin          string `json:"pin"`
	}

	Eth struct {
		Endpoint        string `json:"endpoint,omitempty"`
		PrivateKey      string `json:"private_key,omitempty"`
		FactoryContract string `json:"factory_contract,omitempty"`
	}

	Config struct {
		DB   db.Config `json:"db"`
		Dapp Dapp      `json:"dapp,omitempty"`
		Eth  Eth       `json:"eth,omitempty"`
		Gas  Gas       `json:"gas,omitempty"`
	}

	Gas struct {
		Min              decimal.Decimal `json:"min"`
		Multiplier       decimal.Decimal `json:"multiplier"`
		StrictMultiplier decimal.Decimal `json:"strict_multiplier"`
	}
)
