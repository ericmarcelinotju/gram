package dto

import (
	"encoding/json"
	"time"
)

type JobDto struct {
	Value     interface{}
	CreatedAt time.Time
	Retry     int
}

func CreateJob(value interface{}) JobDto {
	job := JobDto{
		Value:     value,
		CreatedAt: time.Now(),
		Retry:     0,
	}
	return job
}

func (j JobDto) JSON() []byte {
	res, _ := json.Marshal(j)
	return res
}
