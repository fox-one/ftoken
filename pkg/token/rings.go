package token

import (
	"fmt"

	"github.com/fox-one/ftoken/core"
)

func exportRings(asset *core.Asset) *core.TokenItem {
	var supply uint64 = 10000000000
	return &core.TokenItem{
		Name:        fmt.Sprintf("Pando Rings %s", asset.DisplaySymbol),
		Symbol:      fmt.Sprintf("r%s", asset.DisplaySymbol),
		TotalSupply: supply,
	}
}
