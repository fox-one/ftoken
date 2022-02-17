package core

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
)

type (
	TokenItem struct {
		AssetID     string `json:"asset_id,omitempty"`
		AssetKey    string `json:"asset_key"`
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		TotalSupply uint64 `json:"total_supply"`
	}

	TokenItems []*TokenItem
)

func EncodeTokens(tokens TokenItems) ([]byte, error) {
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

func EncodeToken(token TokenItem) ([]byte, error) {
	enc := bytes.NewBuffer(nil)
	enc.WriteByte(byte(len(token.Name)))
	enc.Write([]byte(token.Name))
	enc.WriteByte(byte(len(token.Symbol)))
	enc.Write([]byte(token.Symbol))
	enc.Write(uint64ToByte(token.TotalSupply))
	return enc.Bytes(), nil
}

func DecodeTokens(data []byte) TokenItems {
	var (
		token  *TokenItem
		tokens TokenItems
	)

	for len(data) > 10 {
		if token, data = DecodeToken(data); token != nil {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

func DecodeToken(data []byte) (*TokenItem, []byte) {
	if len(data) <= 10 {
		return nil, nil
	}

	var token TokenItem
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

func (t TokenItem) MarshalBinary() ([]byte, error) {
	return EncodeToken(t)
}

func (t *TokenItem) UnmarshalBinary(data []byte) error {
	token, _ := DecodeToken(data)
	if token == nil {
		return errors.New("unmarshal Token failed")
	}
	*t = *token
	return nil
}

func (t TokenItems) MarshalBinary() ([]byte, error) {
	return EncodeTokens(t)
}

func (t *TokenItems) UnmarshalBinary(data []byte) error {
	tokens := DecodeTokens(data)
	*t = tokens
	return nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (s *TokenItems) Scan(value interface{}) error {
	var d []byte
	switch v := value.(type) {
	case string:
		d = []byte(v)
	case []byte:
		d = v
	}
	var tokens TokenItems
	if err := json.Unmarshal(d, &tokens); err != nil {
		return err
	}
	*s = tokens
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (s TokenItems) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return data, err
}
