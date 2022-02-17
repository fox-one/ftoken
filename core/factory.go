package core

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	TransactionStateNew TransactionState = iota
	TransactionStatePending
	TransactionStateFailed
	TransactionStateSuccess
)

type (
	TransactionState int

	Transaction struct {
		ID        uint64           `sql:"PRIMARY_KEY;" json:"id"`
		CreatedAt time.Time        `json:"created_at"`
		UpdatedAt time.Time        `json:"updated_at"`
		Version   int              `json:"version"`
		TraceID   string           `sql:"size:36;" json:"trace_id,omitempty"`
		Hash      string           `sql:"size:127;" json:"hash,omitempty"`
		Raw       string           `sql:"type:longtext;" json:"raw,omitempty"`
		State     TransactionState `json:"state,omitempty"`
		Tokens    TokenItems       `sql:"type:longtext;" json:"tokens,omitempty"`
		Gas       decimal.Decimal  `sql:"type:decimal(64,8)" json:"gas,omitempty"`
	}

	Factory interface {
		Platform() string
		GasAsset() string
		CreateTransaction(ctx context.Context, tokens TokenItems, trace string) (*Transaction, error)
		SendTransaction(ctx context.Context, tx *Transaction) error
		ReadTransaction(ctx context.Context, hash string) (*Transaction, error)
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Update(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, hash string) (*Transaction, error)
		FindTrace(ctx context.Context, traceID string) ([]*Transaction, error)
	}
)
