package asset

import (
	"context"

	"github.com/fox-one/ftoken/core"
	"github.com/fox-one/pkg/store/db"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Asset{})

		if err := tx.AutoMigrate(core.Asset{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB) core.AssetStore {
	return &assetStore{
		db: db,
	}
}

type assetStore struct {
	db *db.DB
}

func (s *assetStore) Save(ctx context.Context, asset *core.Asset) error {
	return s.db.Update().Model(asset).
		Where("asset_id = ?", asset.AssetID).
		Assign(toUpdateParams(asset)).FirstOrCreate(asset).Error
}

func (s *assetStore) Find(ctx context.Context, ids ...string) ([]*core.Asset, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var assets []*core.Asset
	if err := s.db.View().Where("asset_id IN (?)", ids).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *assetStore) ListAll(ctx context.Context) ([]*core.Asset, error) {
	var assets []*core.Asset
	if err := s.db.View().Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func toUpdateParams(asset *core.Asset) map[string]interface{} {
	params := map[string]interface{}{
		"version": asset.Version + 1,
	}
	if asset.Verified {
		params["verified"] = asset.Verified
	}
	if asset.DisplaySymbol != "" {
		params["display_symbol"] = asset.DisplaySymbol
	}
	return params
}

func (s *assetStore) Update(ctx context.Context, asset *core.Asset) error {
	params := toUpdateParams(asset)
	if tx := s.db.Update().Model(asset).Where("version = ?", asset.Version).Update(params); tx.Error != nil {
		return tx.Error
	} else if tx.RowsAffected == 0 {
		return db.ErrOptimisticLock
	}
	return nil
}
