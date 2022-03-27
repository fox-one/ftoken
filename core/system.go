package core

import (
	"github.com/shopspring/decimal"
)

type (
	Fee struct {
		FeeAssetID string
		FeeAmount  decimal.Decimal
	}

	// System stores system information.
	System struct {
		Version      string
		ClientID     string
		ClientSecret string
		Fees         map[string]*Fee
		WhiteList    map[string]bool
	}
)
