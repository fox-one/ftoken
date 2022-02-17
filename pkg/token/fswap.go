package token

import (
	"fmt"
	"sort"

	"github.com/fox-one/ftoken/core"
)

var (
	Assets = []string{
		"31d2ea9c-95eb-3355-b65b-ba096853bc18", // pUSD
		"4d8c508b-91c5-375b-92b0-ee702ed2dac5", // USDT
		"b91e18ff-a9ae-3dc7-8679-e935d9a4b34b", // USDT@TRON
		"5dac5e28-ad13-31ea-869f-41770dfcee09", // USDT@EOS
		"815b0b1a-2764-3736-8faa-42d694fa620a", // USDT@OMNI
		"9b180ab6-6abe-3dc0-a13f-04169eb34bfa", // USDC
		"0ff3f325-4f34-334d-b6c0-a3bd8850fc06", // JPYC
		"c6d0c728-2624-429b-8e0d-d9d19b6592fa", // BTC
		"43d61dcd-e413-450d-80b8-101d5e903357", // ETH
		"c94ac88f-4671-3976-b60a-09064f1811e8", // XIN
	}
)

func assetIndex(assetID string) int {
	for i, asset := range Assets {
		if asset == assetID {
			return i
		}
	}
	return 1000000
}

func exportFswapLP(assets []*core.Asset) *core.TokenItem {
	sort.Slice(assets, func(i, j int) bool {
		i1 := assetIndex(assets[i].AssetID)
		i2 := assetIndex(assets[j].AssetID)
		if i1 > i2 {
			return true
		} else if i1 < i2 {
			return false
		}

		return assets[i].DisplaySymbol < assets[j].DisplaySymbol
	})

	pair := fmt.Sprintf("%s-%s", assets[0].DisplaySymbol, assets[1].DisplaySymbol)

	var supply uint64 = 10000000000
	return &core.TokenItem{
		Name:        fmt.Sprintf("4swap LP Token %s", pair),
		Symbol:      fmt.Sprintf("s%s", pair),
		TotalSupply: supply,
	}
}
