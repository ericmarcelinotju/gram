package backup

import (
	"context"

	"gitlab.com/firelogik/helios/constant/enums"
	"gitlab.com/firelogik/helios/data/job"
	"gitlab.com/firelogik/helios/domain/model"
	logDomain "gitlab.com/firelogik/helios/domain/module/log"
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
		Title:   "Backup Scheduler Producer Problem",
		Subject: subject,
		Content: content,
		Level:   level,
		Type:    enums.LogTypeSystem,
	})
}
