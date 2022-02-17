package token

import (
	"context"
	"errors"

	"github.com/fox-one/ftoken/core"
)

var (
	ErrInvalidType     = errors.New("invalid token type")
	ErrAssetUnverified = errors.New("asset unverified")
)

func ExportTokenItems(ctx context.Context, assets core.AssetStore, reqs core.TokenRequests) (core.TokenItems, error) {
	tokens := make(core.TokenItems, 0, len(reqs))
	for _, req := range reqs {
		if token, err := ExportTokenItem(ctx, assets, req); err != nil {
			return nil, err
		} else {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}

func ExportTokenItem(ctx context.Context, assets core.AssetStore, req *core.TokenRequest) (*core.TokenItem, error) {
	switch req.Type {
	case core.TokenTypeFswapLP:
		assets, err := assets.Find(ctx, req.Asset1, req.Asset2)
		if err != nil {
			return nil, err
		} else if len(assets) == 2 && assetsVerified(ctx, assets) {
			return exportFswapLP(assets), nil
		}
		return nil, ErrAssetUnverified

	case core.TokenTypeRings:
		assets, err := assets.Find(ctx, req.Asset1)
		if err != nil {
			return nil, err
		} else if len(assets) == 1 && assetsVerified(ctx, assets) {
			return exportRings(assets[0]), nil
		}
		return nil, ErrAssetUnverified

	default:
		return nil, ErrInvalidType
	}
}

func assetsVerified(ctx context.Context, assets []*core.Asset) bool {
	for _, asset := range assets {
		if !asset.Verified {
			return false
		}
	}
	return true
}
