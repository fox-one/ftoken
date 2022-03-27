package config

import (
	"math/big"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/shopspring/decimal"
)

type (
	Config struct {
		DB        db.Config `json:"db"`
		Dapp      Dapp      `json:"dapp,omitempty"`
		Eth       Eth       `json:"eth,omitempty"`
		Fees      []Fee     `json:"fees,omitempty"`
		WhiteList []string  `json:"white_list,omitempty"`
	}

	Fee struct {
		Platform   string          `json:"platform"`
		FeeAssetID string          `json:"fee_asset_id,omitempty"`
		FeeAmount  decimal.Decimal `json:"fee_amount,omitempty"`
	}

	Dapp struct {
		mixin.Keystore
		ClientSecret string `json:"client_secret"`
		Pin          string `json:"pin"`
	}

	Eth struct {
		Endpoint        string   `json:"endpoint,omitempty"`
		PrivateKey      string   `json:"private_key,omitempty"`
		FactoryContract string   `json:"factory_contract,omitempty"`
		MaxGasPrice     *big.Int `json:"max_gas_price,omitempty"`
	}
)
