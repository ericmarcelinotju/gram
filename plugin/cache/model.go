package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(ctx context.Context, key string, data interface{}, expiry time.Duration) (err error)
	Del(ctx context.Context, keys ...string) (err error)
	Get(ctx context.Context, key string, data interface{}) (err error)
	HashGet(ctx context.Context, key, field string, data interface{}) (err error)
	HashGetAll(ctx context.Context, key string) (result map[string]string, err error)
	HashSet(ctx context.Context, pkey, field string, data interface{}) (err error)
	HashDel(ctx context.Context, key, field string) (err error)

	Client() *redis.Client
}
