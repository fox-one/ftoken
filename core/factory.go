package core

import (
	"bytes"
	"context"
	"encoding/binary"
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

	Token struct {
		Name        string `gorm:"size:255;" json:"name"`
		Symbol      string `gorm:"size:255;" json:"symbol"`
		TotalSupply uint64 `json:"total_supply"`
		AssetKey    string `gorm:"size:255;" json:"asset_key"`
		AssetID     string `gorm:"size:36;" json:"asset_id"`
	}

	Tokens []*Token

	Transaction struct {
		ID        uint64           `gorm:"PRIMARY_KEY;" json:"id"`
		CreatedAt time.Time        `json:"created_at"`
		UpdatedAt time.Time        `json:"updated_at"`
		TraceID   string           `gorm:"size:36;" json:"trace_id,omitempty"`
		Hash      string           `json:"hash,omitempty"`
		Raw       string           `gorm:"type:longtext;" json:"raw,omitempty"`
		State     TransactionState `json:"state,omitempty"`
		Tokens    Tokens           `json:"tokens,omitempty"`
		Gas       decimal.Decimal  `json:"gas,omitempty"`
	}

	Factory interface {
		Platform() string
		GasAsset() string
		CreateTransaction(ctx context.Context, tokens []*Token, receiver string) (*Transaction, error)
		SendTransaction(ctx context.Context, tx *Transaction) error
		ReadTransaction(ctx context.Context, hash string) (*Transaction, error)
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Update(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, hash string) (*Transaction, error)
	}
)

func EncodeTokens(tokens Tokens) ([]byte, error) {
	enc := bytes.NewBuffer(nil)
	for _, token := range tokens {
		enc.WriteByte(byte(len(token.Name)))
		enc.Write([]byte(token.Name))
		enc.WriteByte(byte(len(token.Symbol)))
		enc.Write([]byte(token.Symbol))
		enc.Write(uint64ToByte(token.TotalSupply))
	}
	return enc.Bytes(), nil
}

func uint64ToByte(d uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, d)
	return b
}
