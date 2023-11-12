package job

import (
	"encoding/json"
	"time"
)

type Job struct {
	Value     interface{}
	CreatedAt time.Time
	Retry     int
}

func CreateJob(value interface{}) Job {
	job := Job{
		Value:     value,
		CreatedAt: time.Now(),
		Retry:     0,
	}
	return job
}

func (j Job) JSON() []byte {
	res, _ := json.Marshal(j)
	return res
}
