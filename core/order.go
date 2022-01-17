package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	OrderStatePending OrderState = iota
	OrderStatePaid
	OrderStateFailed
	OrderStateDone
)

type (
	OrderState int

	Order struct {
		ID          uint64          `gorm:"PRIMARY_KEY;" json:"id"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
		Version     int             `json:"version"`
		ValidBefore time.Time       `json:"valid_before"`
		TraceID     string          `sql:"size:36;" json:"trace_id"`
		State       OrderState      `json:"state"`
		UserID      string          `gorm:"size:36;" json:"user_id"`
		FeeAsset    string          `gorm:"size:36;" json:"fee_asset"`
		FeeAmount   decimal.Decimal `sql:"type:decimal(64,8)" json:"fee_amount"`
		Platform    string          `gorm:"size:255;" json:"platform"`
		Tokens      Tokens          `gorm:"type:longtext;" json:"tokens"`
		Receiver    string          `gorm:"size:255;" json:"receiver"`
		Transaction string          `gorm:"size:255;" json:"transaction"`
	}

	OrderStore interface {
		Create(ctx context.Context, order *Order) error
		Update(ctx context.Context, order *Order) error
		Find(ctx context.Context, traceID string) (*Order, error)
		List(ctx context.Context, from uint64, limit int) ([]*Order, error)
	}
)
