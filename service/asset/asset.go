package asset

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/mixin-sdk-go"
)

func New(c *mixin.Client) core.AssetService {
	return &assetService{c: c}
}

type assetService struct {
	c *mixin.Client
}

func (s *assetService) Find(ctx context.Context, assetID string) (*core.Asset, error) {
	asset, err := s.c.ReadAsset(ctx, assetID)
	if err != nil {
		if mixin.IsErrorCodes(err, 10002) {
			err = core.ErrAssetNotExist
		}

		return nil, err
	}

	return convertAsset(asset), nil
}

func convertAsset(asset *mixin.Asset) *core.Asset {
	return &core.Asset{
		AssetID: asset.AssetID,
		Name:    asset.Name,
		Symbol:  asset.Symbol,
		ChainID: asset.ChainID,
		Logo:    asset.IconURL,
	}
}
