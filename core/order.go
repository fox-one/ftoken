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
	OrderStatePending OrderState = iota
	OrderStatePaid
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
		ID        uint64          `gorm:"PRIMARY_KEY;" json:"id"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
		Version   int             `json:"version"`
		TraceID   string          `sql:"size:36;" json:"trace_id"`
		State     OrderState      `json:"state"`
		UserID    string          `gorm:"size:36;" json:"user_id"`
		FeeAsset  string          `gorm:"size:36;" json:"fee_asset"`
		FeeAmount decimal.Decimal `sql:"type:decimal(64,8)" json:"fee_amount"`
		Platform  string          `gorm:"size:255;" json:"platform"`
		Tokens    Tokens          `gorm:"type:longtext;" json:"tokens"`
		Result    Tokens          `gorm:"type:longtext;" json:"result"`
		Receiver  *Address        `gorm:"size:255;" json:"receiver"`
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
