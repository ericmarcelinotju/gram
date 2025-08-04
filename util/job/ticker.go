package job

import (
	"time"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour

const HOUR_TO_TICK int = 24
const MINUTE_TO_TICK int = 00
const SECOND_TO_TICK int = 00

type JobTicker struct {
	Timer *time.Timer
}

func (t *JobTicker) UpdateTimer() {
	nextTick := time.Date(time.Now().Year(), time.Now().Month(),
		time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
	if !nextTick.After(time.Now()) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	diff := nextTick.Sub(time.Now())
	if t.Timer == nil {
		t.Timer = time.NewTimer(diff)
	} else {
		t.Timer.Reset(diff)
	}
}
