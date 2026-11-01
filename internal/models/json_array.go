package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// StringArray is a helper type to store []string as JSON in PostgreSQL
type StringArray []string

// Value converts StringArray to JSON byte slice for database driver
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

// Scan converts database value to StringArray
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("type assertion to []byte/string failed")
	}

	if len(bytes) == 0 {
		*a = []string{}
		return nil
	}

	return json.Unmarshal(bytes, a)
}
