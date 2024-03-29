package core

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

var (
	ErrInvalidTrace = errors.New("invalid trace")
)

const (
	TransferPriorityHigh TransferPriority = iota
	TransferPriorityNormal
	TransferPriorityLow
)

type (
	Snapshot struct {
		CreatedAt       time.Time       `json:"created_at,omitempty"`
		SnapshotID      string          `json:"snapshot_id,omitempty"`
		UserID          string          `json:"user_id,omitempty"`
		OpponentID      string          `json:"opponent_id,omitempty"`
		TraceID         string          `json:"trace_id,omitempty"`
		AssetID         string          `json:"asset_id,omitempty"`
		Source          string          `json:"source,omitempty"`
		Amount          decimal.Decimal `json:"amount,omitempty"`
		Memo            string          `json:"memo,omitempty"`
		TransactionHash string          `json:"transaction_hash,omitempty"`
	}

	TransferPriority int

	Transfer struct {
		ID        int64            `sql:"PRIMARY_KEY" json:"id,omitempty"`
		Priority  TransferPriority `json:"priority"`
		CreatedAt time.Time        `json:"created_at,omitempty"`
		UpdatedAt time.Time        `json:"updated_at,omitempty"`
		TraceID   string           `sql:"type:char(36)" json:"trace_id,omitempty"`
		AssetID   string           `sql:"type:char(36)" json:"asset_id,omitempty"`
		Amount    decimal.Decimal  `sql:"type:decimal(64,8)" json:"amount,omitempty"`
		Memo      string           `sql:"size:200" json:"memo,omitempty"`
		Threshold uint8            `json:"threshold,omitempty"`
		Opponents pq.StringArray   `sql:"type:varchar(1024)" json:"opponents,omitempty"`
	}

	WalletStore interface {
		ListTransfers(ctx context.Context, limit int) ([]*Transfer, error)
		CreateTransfers(ctx context.Context, transfers []*Transfer) error
		ExpireTransfers(ctx context.Context, transfers []*Transfer) error
		CountTransfers(ctx context.Context) (int, error)
	}

	WalletService interface {
		ListSnapshots(ctx context.Context, offset time.Time, limit int) ([]*Snapshot, error)
		Transfer(ctx context.Context, transfer *Transfer) error
	}
)
