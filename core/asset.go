package core

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrAssetNotExist = errors.New("asset not exist")
)

type (
	Asset struct {
		ID             uint64          `sql:"PRIMARY_KEY;" json:"id"`
		AssetID        string          `sql:"size:36;" json:"asset_id,omitempty"`
		CreatedAt      time.Time       `json:"created_at,omitempty"`
		UpdatedAt      time.Time       `json:"updated_at,omitempty"`
		Version        int64           `sql:"not null" json:"version,omitempty"`
		Verified       bool            `json:"verified"`
		Name           string          `sql:"size:64" json:"name,omitempty"`
		Symbol         string          `sql:"size:32" json:"symbol,omitempty"`
		Logo           string          `sql:"size:256" json:"logo,omitempty"`
		DisplaySymbol  string          `sql:"size:32;default:null;" json:"display_symbol,omitempty"`
		ChainID        string          `sql:"size:36" json:"chain_id,omitempty"`
		Price          decimal.Decimal `sql:"type:decimal(32,8);default:0" json:"price,omitempty"`
		PriceUpdatedAt *time.Time      `sql:"default:null" json:"price_updated_at,omitempty"`
	}

	// AssetStore defines operations for working with assets on db.
	AssetStore interface {
		Save(ctx context.Context, asset *Asset) error
		Find(ctx context.Context, assetIDs ...string) ([]*Asset, error)
		ListAll(ctx context.Context) ([]*Asset, error)
	}

	// AssetService provides access to assets information
	// in the remote system like mixin network.
	AssetService interface {
		Find(ctx context.Context, assetID string) (*Asset, error)
	}
)
