package enums

import (
	"encoding/json"
)

type LogLevel string

const (
	DefaultLogLevel LogLevel = ""
	LogLevelInfo    LogLevel = "info"
	LogLevelSuccess LogLevel = "success"
	LogLevelWarning LogLevel = "warning"
	LogLevelDanger  LogLevel = "danger"
)

var logLevelIds = map[string]LogLevel{
	"info":    LogLevelInfo,
	"success": LogLevelSuccess,
	"warning": LogLevelWarning,
	"danger":  LogLevelDanger,
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *LogLevel) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*s = logLevelIds[j]
	return nil
}
