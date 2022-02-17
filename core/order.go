package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	OrderStateNew OrderState = iota
	OrderStatePaid
	OrderStateProcessing
	OrderStateFailed
	OrderStateDone
)

type (
	OrderState int

	Order struct {
		ID            uint64          `sql:"PRIMARY_KEY;" json:"id"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
		Version       int             `json:"version,omitempty"`
		TraceID       string          `sql:"size:36;" json:"trace_id,omitempty"`
		State         OrderState      `json:"state"`
		UserID        string          `sql:"size:36;" json:"user_id,omitempty"`
		FeeAsset      string          `sql:"size:36;" json:"fee_asset,omitempty"`
		FeeAmount     decimal.Decimal `sql:"type:decimal(64,8)" json:"fee_amount,omitempty"`
		GasUsage      decimal.Decimal `sql:"type:decimal(64,8)" json:"gas_usage,omitempty"`
		Platform      string          `sql:"size:255;" json:"platform,omitempty"`
		TokenRequests TokenItems      `sql:"type:longtext;" json:"token_requests,omitempty"`
		Tokens        TokenItems      `sql:"type:longtext;" json:"tokens,omitempty"`
		Transaction   string          `sql:"size:128;" json:"transaction,omitempty"`
	}

	OrderStore interface {
		Create(ctx context.Context, order *Order) error
		Update(ctx context.Context, order *Order) error
		Find(ctx context.Context, traceID string) (*Order, error)
		List(ctx context.Context, state OrderState, limit int) ([]*Order, error)
	}
)
