package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gitlab.com/firelogik/helios/config"
)

func ConnectRedis(configuration *config.Cache) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configuration.Host, configuration.Port),
		Password: configuration.Password,
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return NewRedisCache(client), nil
}
