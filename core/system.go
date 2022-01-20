package core

import "github.com/shopspring/decimal"

type (
	Gas struct {
		Min              decimal.Decimal
		Multiplier       decimal.Decimal
		StrictMultiplier decimal.Decimal
	}

	// System stores system information.
	System struct {
		Version      string
		ClientID     string
		ClientSecret string
		Gas          Gas
		Addresses    map[string]*Address
	}
)
