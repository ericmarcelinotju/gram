package example

import (
	"context"

	"github.com/ericmarcelinotju/gram/constant/enums"
	"github.com/ericmarcelinotju/gram/dto"
)

// NewProducerFactory create and returns a factory to create routes for the panelment
func NewProducerFactory() func(ctx context.Context) ([]dto.JobDto, error) {
	backupProducer := func(ctx context.Context) (jobs []dto.JobDto, err error) {
		// Create jobs
		return jobs, nil
	}
	return backupProducer
}

func CreateProduceLog(ctx context.Context, subject, content string, level enums.LogLevel) {
	// Log here
}
