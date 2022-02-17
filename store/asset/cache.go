package asset

import (
	"context"
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

func Cache(store core.AssetStore, exp time.Duration) core.AssetStore {
	return &cacheAssetStore{
		AssetStore: store,
		cache:      cache.New(exp, cache.NoExpiration),
		sf:         &singleflight.Group{},
	}
}

type cacheAssetStore struct {
	core.AssetStore
	cache *cache.Cache
	sf    *singleflight.Group
}

func (s *cacheAssetStore) Save(ctx context.Context, asset *core.Asset) error {
	s.cache.Delete(asset.AssetID)
	if err := s.AssetStore.Save(ctx, asset); err != nil {
		return err
	}

	s.cache.Delete(asset.AssetID)
	return nil
}

func (s *cacheAssetStore) Find(ctx context.Context, ids ...string) ([]*core.Asset, error) {
	if assets, ok := s.itemsFromCache(ctx, ids...); ok {
		return assets, nil
	}

	assets, err := s.AssetStore.Find(ctx, ids...)
	if err != nil {
		return nil, err
	}

	for _, asset := range assets {
		s.cache.SetDefault(asset.AssetID, asset)
	}

	return assets, nil
}

func (s *cacheAssetStore) itemsFromCache(ctx context.Context, ids ...string) ([]*core.Asset, bool) {
	var assets = make([]*core.Asset, 0, len(ids))
	for i, id := range ids {
		v, ok := s.cache.Get(id)
		if !ok {
			return nil, false
		}
		assets[i] = v.(*core.Asset)
	}
	return assets, true
}
