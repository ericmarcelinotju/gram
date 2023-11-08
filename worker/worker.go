package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/ericmarcelinotju/gram/constant/enums"
	"github.com/ericmarcelinotju/gram/data/job"
	"github.com/ericmarcelinotju/gram/domain/model"
	logDomain "github.com/ericmarcelinotju/gram/domain/module/log"
)

type Worker struct {
	ctx       context.Context
	scheduler *job.Scheduler
	queue     *job.Queue

	logSvc logDomain.Service
}

// StartWorker start worker
func NewWorker(
	ctx context.Context,
	scheduler *job.Scheduler,
	queue *job.Queue,

	logSvc logDomain.Service,
) (*Worker, error) {
	return &Worker{
		ctx:       ctx,
		scheduler: scheduler,
		queue:     queue,

		logSvc: logSvc,
	}, nil
}

// StartWorker start worker
func (w *Worker) Start() error {
	var err error

	err = w.scheduler.SetScheduleFunc(w.OnSchedule)
	if err != nil {
		return err
	}

	err = w.queue.StartConsuming()
	if err != nil {
		if err != rmq.ErrorAlreadyConsuming {
			return err
		}
	}

	consumerFunc := NewConsumerFactory(w.ctx, w.logSvc)
	err = w.queue.AddConsumer(consumerFunc)
	if err != nil {
		return err
	}

	// w.CreateJobs()

	return w.scheduler.Start()
}

func (w *Worker) Stop() error {
	w.scheduler.Stop()
	return nil
}

func (w *Worker) CreateJobs() error {
	backupJobs, err := NewProducerFactory(w.logSvc)(w.ctx)
	if err != nil {
		return err
	}
	for _, job := range backupJobs {
		err := w.queue.Client.PublishBytes(job.JSON())
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) OnSchedule() {
	var log model.Log = model.Log{
		Title:   "Backup Scheduler Started",
		Subject: "Scheduler for backup started with no problem",
		Content: fmt.Sprintf(
			`<b>Scheduler started at</b> : %s`,
			time.Now().Format("02-01-2006 15:04:05"),
		),
		Type:  enums.LogTypeSystem,
		Level: enums.LogLevelInfo,
	}
	err := w.CreateJobs()
	if err != nil {
		log = model.Log{
			Title:   "Backup Scheduler Error",
			Subject: "Error detected when creating jobs for backups",
			Content: fmt.Sprintf(
				`<b>Scheduler started at</b> : %s<br><b>Error</b>: %s`,
				time.Now().Format("02-01-2006 15:04:05"),
				err.Error(),
			),
			Type:  enums.LogTypeSystem,
			Level: enums.LogLevelDanger,
		}
	}
	w.logSvc.CreateLog(w.ctx, &log)
}
