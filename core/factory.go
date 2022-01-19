package core

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
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
		Name        string `gorm:"size:255;" json:"name,omitempty"`
		Symbol      string `gorm:"size:255;" json:"symbol,omitempty"`
		TotalSupply uint64 `json:"total_supply,omitempty"`
		AssetKey    string `gorm:"size:255;" json:"asset_key,omitempty"`
		AssetID     string `gorm:"size:36;" json:"asset_id,omitempty"`
	}

	Tokens []*Token

	Transaction struct {
		ID        uint64           `sql:"PRIMARY_KEY;" json:"id"`
		CreatedAt time.Time        `json:"created_at"`
		UpdatedAt time.Time        `json:"updated_at"`
		Version   int              `json:"version"`
		TraceID   string           `sql:"size:36;" json:"trace_id,omitempty"`
		Hash      string           `json:"hash,omitempty"`
		Raw       string           `sql:"type:longtext;" json:"raw,omitempty"`
		State     TransactionState `json:"state,omitempty"`
		Tokens    Tokens           `sql:"type:longtext;" json:"tokens,omitempty"`
		Gas       decimal.Decimal  `sql:"type:decimal(64,8)" json:"gas,omitempty"`
	}

	Factory interface {
		Platform() string
		GasAsset() string
		CreateTransaction(ctx context.Context, tokens []*Token, receiver *Address) (*Transaction, error)
		SendTransaction(ctx context.Context, tx *Transaction) error
		ReadTransaction(ctx context.Context, hash string) (*Transaction, error)
	}

	TransactionStore interface {
		Create(ctx context.Context, tx *Transaction) error
		Update(ctx context.Context, tx *Transaction) error
		Find(ctx context.Context, traceID string) (*Transaction, error)
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

func EncodeToken(token Token) ([]byte, error) {
	enc := bytes.NewBuffer(nil)
	enc.WriteByte(byte(len(token.Name)))
	enc.Write([]byte(token.Name))
	enc.WriteByte(byte(len(token.Symbol)))
	enc.Write([]byte(token.Symbol))
	enc.Write(uint64ToByte(token.TotalSupply))
	return enc.Bytes(), nil
}

func DecodeTokens(data []byte) Tokens {
	var (
		token  *Token
		tokens Tokens
	)

	for len(data) > 10 {
		if token, data = DecodeToken(data); token != nil {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func DecodeToken(data []byte) (*Token, []byte) {
	if len(data) <= 10 {
		return nil, nil
	}

	var token Token
	offset := 0
	if size := int(data[offset]); len(data) > 10+size {
		offset++
		token.Name = string(data[offset : offset+size])
		offset += size
	} else {
		return nil, nil
	}

	if size := int(data[offset]); len(data) >= offset+size+9 {
		offset++
		token.Symbol = string(data[offset : offset+size])
		offset += size
	} else {
		return nil, nil
	}

	token.TotalSupply = binary.BigEndian.Uint64(data[offset : offset+8])
	offset += 8
	data = data[offset:]
	if token.TotalSupply > 0 {
		return &token, data
	}
	return nil, data
}

func uint64ToByte(d uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, d)
	return b
}

func (t Token) MarshalBinary() ([]byte, error) {
	return EncodeToken(t)
}

func (t *Token) UnmarshalBinary(data []byte) error {
	token, _ := DecodeToken(data)
	if token == nil {
		return errors.New("unmarshal Token failed")
	}
	*t = *token
	return nil
}

func (t Tokens) MarshalBinary() ([]byte, error) {
	return EncodeTokens(t)
}

func (t *Tokens) UnmarshalBinary(data []byte) error {
	tokens := DecodeTokens(data)
	*t = tokens
	return nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (s *Tokens) Scan(value interface{}) error {
	var d []byte
	switch v := value.(type) {
	case string:
		d = []byte(v)
	case []byte:
		d = v
	}
	var tokens Tokens
	if err := json.Unmarshal(d, &tokens); err != nil {
		return err
	}
	*s = tokens
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (s Tokens) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return data, err
}
