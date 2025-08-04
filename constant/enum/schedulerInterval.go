package enums

import (
	"encoding/json"
)

type SchedulerInterval string

const (
	DefaultSchedulerInterval  SchedulerInterval = ""
	SchedulerIntervalMonthly  SchedulerInterval = "monthly"
	SchedulerIntervalDaily    SchedulerInterval = "daily"
	SchedulerIntervalMinutely SchedulerInterval = "minutely"
)

var schedulerIntervalIds = map[string]SchedulerInterval{
	"monthly":  SchedulerIntervalMonthly,
	"daily":    SchedulerIntervalDaily,
	"minutely": SchedulerIntervalMinutely,
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *SchedulerInterval) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*s = schedulerIntervalIds[j]
	return nil
}
