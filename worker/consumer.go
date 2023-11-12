package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ericmarcelinotju/gram/constant/enums"

	"github.com/adjust/rmq/v4"
	"github.com/ericmarcelinotju/gram/data/job"
)

type Consumer struct {
	reportBatchSize int

	name   string
	count  int
	before time.Time

	ctx context.Context
}

func NewJobProcessor() func(ctx context.Context, job job.Job) error {
	return func(ctx context.Context, currJob job.Job) error {
		// Process job

		return nil
	}
}

// NewConsumerFactory create and returns a factory to create billing consumer for scheduler
func NewConsumerFactory(
	ctx context.Context,
) func(tag int, reportBatchSize int) rmq.Consumer {
	consumerFactory := func(tag int, reportBatchSize int) rmq.Consumer {
		return &Consumer{
			reportBatchSize: reportBatchSize,

			name:   fmt.Sprintf("consumer-%d", tag),
			count:  0,
			before: time.Now(),

			ctx: ctx,
		}
	}
	return consumerFactory
}

func (c *Consumer) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()

	c.count++
	if c.count%c.reportBatchSize == 0 {
		duration := time.Since(c.before)
		c.before = time.Now()
		perSecond := time.Second / (duration / time.Duration(c.reportBatchSize))
		CreateConsumeLog(c.ctx, "Report", fmt.Sprintf("%s consumed %d %s %d", c.name, c.count, payload, perSecond), enums.LogLevelInfo)
	}

	var job job.Job
	if err := json.Unmarshal([]byte(payload), &job); err != nil {
		// handle json error
		CreateConsumeLog(c.ctx, "Format job error", err.Error(), enums.LogLevelDanger)
		if err := delivery.Reject(); err != nil {
			// handle reject error
			CreateConsumeLog(c.ctx, "Reject job error", err.Error(), enums.LogLevelWarning)
		}
		return
	}

	err := NewJobProcessor()(c.ctx, job)
	if err != nil {
		// handle error
		CreateConsumeLog(c.ctx, "Process job error", err.Error(), enums.LogLevelDanger)
		if err := delivery.Reject(); err != nil {
			// handle reject error
			CreateConsumeLog(c.ctx, "Reject job error", err.Error(), enums.LogLevelWarning)
		}
		return
	}

	CreateConsumeLog(c.ctx, "Performing task", fmt.Sprintf("Process recording: %v", job.Value), enums.LogLevelInfo)
	if err := delivery.Ack(); err != nil {
		// handle ack error
		CreateConsumeLog(c.ctx, "Acknowledge job error", err.Error(), enums.LogLevelWarning)
	}
}

func CreateConsumeLog(ctx context.Context, subject, content string, level enums.LogLevel) {
	// Log here
}
