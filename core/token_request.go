package core

import (
	"database/sql/driver"
	"encoding/json"
)

const (
	TokenTypeFswapLP TokenType = iota
	TokenTypeRings
)

type (
	TokenType int

	TokenRequest struct {
		Type   TokenType `json:"type"`
		Asset1 string    `json:"asset1,omitempty"`
		Asset2 string    `json:"asset2,omitempty"`
	}

	TokenRequests []*TokenRequest
)

// Scan implements the sql.Scanner interface for database deserialization.
func (s *TokenRequests) Scan(value interface{}) error {
	var d []byte
	switch v := value.(type) {
	case string:
		d = []byte(v)
	case []byte:
		d = v
	}
	var tokens TokenRequests
	if err := json.Unmarshal(d, &tokens); err != nil {
		return err
	}
	*s = tokens
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (s TokenRequests) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return data, err
}
