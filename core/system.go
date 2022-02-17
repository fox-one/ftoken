package core

import (
	"github.com/fox-one/pkg/number"
	"github.com/shopspring/decimal"
)

type (
	Gas struct {
		Mins             number.Values
		Multiplier       decimal.Decimal
		StrictMultiplier decimal.Decimal
	}

	// System stores system information.
	System struct {
		Version      string
		ClientID     string
		ClientSecret string
		Gas          Gas
	}
)
