package views

import (
	"time"

	"github.com/fox-one/ftoken/core"
	"github.com/shopspring/decimal"
)

type (
	Order struct {
		ID            uint64           `json:"id,omitempty"`
		CreatedAt     *time.Time       `json:"created_at,omitempty"`
		UpdatedAt     *time.Time       `json:"updated_at,omitempty"`
		TraceID       string           `json:"trace_id,omitempty"`
		State         core.OrderState  `json:"state,omitempty"`
		UserID        string           `json:"user_id,omitempty"`
		FeeAsset      string           `json:"fee_asset,omitempty"`
		FeeAmount     decimal.Decimal  `json:"fee_amount,omitempty"`
		Platform      string           `json:"platform,omitempty"`
		TokenRequests core.TokenItems  `json:"tokens,omitempty"`
		Tokens        *core.TokenItems `json:"result,omitempty"`
	}
)

func OrderView(o core.Order) Order {
	order := Order{
		ID:            o.ID,
		TraceID:       o.TraceID,
		State:         o.State,
		FeeAsset:      o.FeeAsset,
		FeeAmount:     o.FeeAmount,
		Platform:      o.Platform,
		TokenRequests: o.TokenRequests,
	}

	if o.CreatedAt.IsZero() {
		order.CreatedAt = &o.CreatedAt
	}
	if o.UpdatedAt.IsZero() {
		order.UpdatedAt = &o.UpdatedAt
	}
	if len(o.Tokens) > 0 {
		order.Tokens = &o.Tokens
	}

	return order
}
