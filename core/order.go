package core

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/fox-one/ftoken/pkg/mtg"
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

	Address struct {
		Destination string `json:"destination,omitempty"`
		Tag         string `json:"tag,omitempty"`
	}

	Order struct {
		ID          uint64          `sql:"PRIMARY_KEY;" json:"id"`
		CreatedAt   time.Time       `json:"created_at"`
		UpdatedAt   time.Time       `json:"updated_at"`
		Version     int             `json:"version,omitempty"`
		TraceID     string          `sql:"size:36;" json:"trace_id,omitempty"`
		State       OrderState      `json:"state"`
		UserID      string          `sql:"size:36;" json:"user_id,omitempty"`
		FeeAsset    string          `sql:"size:36;" json:"fee_asset,omitempty"`
		FeeAmount   decimal.Decimal `sql:"type:decimal(64,8)" json:"fee_amount,omitempty"`
		GasUsage    decimal.Decimal `sql:"type:decimal(64,8)" json:"gas_usage,omitempty"`
		Platform    string          `sql:"size:255;" json:"platform,omitempty"`
		Tokens      Tokens          `sql:"type:longtext;" json:"tokens,omitempty"`
		Result      Tokens          `sql:"type:longtext;" json:"result,omitempty"`
		Receiver    *Address        `sql:"size:255;" json:"receiver,omitempty"`
		Transaction string          `sql:"size:128;" json:"transaction,omitempty"`
	}

	OrderStore interface {
		Create(ctx context.Context, order *Order) error
		Update(ctx context.Context, order *Order) error
		Find(ctx context.Context, traceID string) (*Order, error)
		List(ctx context.Context, state OrderState, limit int) ([]*Order, error)
	}
)

func (a Address) MarshalBinary() ([]byte, error) {
	return mtg.Encode(a.Destination, a.Tag)
}

func (a *Address) UnmarshalBinary(data []byte) error {
	var (
		destination string
		tag         string
	)
	if _, err := mtg.Scan(data, &destination, &tag); err != nil {
		return err
	}
	a.Destination = destination
	a.Tag = tag
	return nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (a *Address) Scan(value interface{}) error {
	var d []byte
	switch v := value.(type) {
	case string:
		d = []byte(v)
	case []byte:
		d = v
	}
	var address Address
	if err := json.Unmarshal(d, &address); err != nil {
		return err
	}
	*a = address
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (a *Address) Value() (driver.Value, error) {
	return json.Marshal(a)
}
