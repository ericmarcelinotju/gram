package job

import (
	"fmt"
	"log"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/ericmarcelinotju/gram/config"
	"github.com/go-redis/redis/v8"
)

type Queue struct {
	Client     rmq.Queue
	Connection rmq.Connection
	Config     config.Queue
}

func ConnectQueue(config *config.Queue, client *redis.Client) (*Queue, error) {
	errChan := make(chan error, 10)
	connection, err := rmq.OpenConnectionWithRedisClient("queue", client, errChan)
	if err != nil {
		return nil, err
	}
	queue, err := connection.OpenQueue(config.Name)
	if err != nil {
		return nil, err
	}

	cleaner := rmq.NewCleaner(connection)

	go (func() {
		for range time.Tick(time.Minute) {
			_, err := cleaner.Clean()
			if err != nil {
				log.Printf("[QUEUE] failed to clean: %s", err)
				continue
			}
			// log.Printf("[QUEUE] cleaned %d", returned)
		}
	})()

	return &Queue{
		Client:     queue,
		Connection: connection,
		Config:     *config,
	}, nil
}

func (q *Queue) StartConsuming() error {
	return q.Client.StartConsuming(q.Config.PrefetchLimit, q.Config.PollDuration)
}

func (q *Queue) AddConsumer(consumerFactory func(int, int) rmq.Consumer) error {
	for i := 0; i < q.Config.Number; i++ {
		name := fmt.Sprintf("consumer %d", i)
		if _, err := q.Client.AddConsumer(name, consumerFactory(i, q.Config.ReportBatchSize)); err != nil {
			return err
		}
	}
	return nil
}
