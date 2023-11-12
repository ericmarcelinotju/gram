package worker

import (
	"context"

	"github.com/adjust/rmq/v4"
	"github.com/ericmarcelinotju/gram/data/job"
)

type Worker struct {
	ctx       context.Context
	scheduler *job.Scheduler
	queue     *job.Queue
}

// StartWorker start worker
func NewWorker(
	ctx context.Context,
	scheduler *job.Scheduler,
	queue *job.Queue,
) (*Worker, error) {
	return &Worker{
		ctx:       ctx,
		scheduler: scheduler,
		queue:     queue,
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

	consumerFunc := NewConsumerFactory(w.ctx)
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
	backupJobs, err := NewProducerFactory()(w.ctx)
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
	err := w.CreateJobs()
	if err != nil {
		// Log here
	}
	// Log here
}
