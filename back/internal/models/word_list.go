package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type WordList []string

func (s WordList) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *WordList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}
