package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSON string

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	val, ok := value.(json.RawMessage)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	jsonStr, err := json.Marshal(val)
	*j = JSON(jsonStr)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	value, err := json.RawMessage(j).MarshalJSON()
	return value, err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) RawMessage() json.RawMessage {
	if len(j) == 0 {
		return json.RawMessage("{}")
	}
	return json.RawMessage(j)
}
