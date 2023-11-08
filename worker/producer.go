package worker

import (
	"context"

	"github.com/ericmarcelinotju/gram/constant/enums"
	"github.com/ericmarcelinotju/gram/data/job"
	"github.com/ericmarcelinotju/gram/domain/model"
	logDomain "github.com/ericmarcelinotju/gram/domain/module/log"
)

// NewProducerFactory create and returns a factory to create routes for the panelment
func NewProducerFactory(logSvc logDomain.Service) func(ctx context.Context) ([]job.Job, error) {
	backupProducer := func(ctx context.Context) (jobs []job.Job, err error) {
		// Create jobs
		return jobs, nil
	}
	return backupProducer
}

func CreateProduceLog(logSvc logDomain.Service, ctx context.Context, subject, content string, level enums.LogLevel) {
	logSvc.CreateLog(ctx, &model.Log{
		Title:   "Scheduler Producer Problem",
		Subject: subject,
		Content: content,
		Level:   level,
		Type:    enums.LogTypeSystem,
	})
}
