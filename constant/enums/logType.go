package enums

import (
	"encoding/json"
)

type LogType string

const (
	DefaultLogType LogType = ""
	LogTypeEvent   LogType = "event"
	LogTypeSystem  LogType = "system"
)

var logTypeIds = map[string]LogType{
	"event":  LogTypeEvent,
	"system": LogTypeSystem,
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *LogType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*s = logTypeIds[j]
	return nil
}
