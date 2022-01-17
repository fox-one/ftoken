package config

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
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
	}
)
