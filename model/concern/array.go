package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		return nil // case when value from the db was NULL
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to cast value to string: %v", value)
	}
	if len(s) > 0 {
		*a = strings.Split(s, ";")
	}
	return nil
}

// Value return json value, implement driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	return strings.Join(a, ";"), nil
}

type IntArray []int

func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		return nil // case when value from the db was NULL
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to cast value to string: %v", value)
	}
	if len(s) > 0 {
		strings := strings.Split(s, ";")
		intArr := make([]int, len(strings))
		for i, str := range strings {
			intItem, err := strconv.Atoi(str)
			if err != nil {
				return err
			}
			intArr[i] = intItem
		}
		*a = intArr
	}
	return nil
}

// Value return json value, implement driver.Valuer interface
func (a IntArray) Value() (driver.Value, error) {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), ";"), "[]"), nil
}
