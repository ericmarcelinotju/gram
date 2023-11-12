package worker

import (
	"context"

	"github.com/ericmarcelinotju/gram/constant/enums"
	"github.com/ericmarcelinotju/gram/data/job"
)

// NewProducerFactory create and returns a factory to create routes for the panelment
func NewProducerFactory() func(ctx context.Context) ([]job.Job, error) {
	backupProducer := func(ctx context.Context) (jobs []job.Job, err error) {
		// Create jobs
		return jobs, nil
	}
	return backupProducer
}

func CreateProduceLog(ctx context.Context, subject, content string, level enums.LogLevel) {
	// Log here
}
