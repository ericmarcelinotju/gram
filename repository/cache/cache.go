package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

//CONSTRUCTOR STRUCT
type redisCache struct {
	cache *redis.Client
}

//CONSTRUCTOR FUNCTION FOR USER REPOSITORY
func NewRedisCache(cache *redis.Client) Cache {
	return &redisCache{cache: cache}
}

func (r *redisCache) Client() *redis.Client {
	return r.cache
}

func (r *redisCache) Set(ctx context.Context, key string, data interface{}, expiry time.Duration) (err error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = r.cache.Set(ctx, key, string(jsonBytes), expiry).Err()
	return
}

func (r *redisCache) Del(ctx context.Context, key ...string) (err error) {
	err = r.cache.Del(ctx, key...).Err()
	return
}

func (r *redisCache) Get(ctx context.Context, key string, data interface{}) (err error) {
	res, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(res), &data)
	return
}

func (r *redisCache) HashGet(ctx context.Context, key, field string, data interface{}) (err error) {
	result, err := r.cache.HGet(ctx, key, field).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(result), data)
	return
}

func (r *redisCache) HashGetAll(ctx context.Context, key string) (result map[string]string, err error) {
	result, err = r.cache.HGetAll(ctx, key).Result()
	return
}

func (r *redisCache) HashSet(ctx context.Context, key, field string, data interface{}) (err error) {
	execute := r.cache.HSet(ctx, key, field, data)
	return execute.Err()
}

func (r *redisCache) HashDel(ctx context.Context, key, field string) (err error) {
	_, err = r.cache.HDel(ctx, key, field).Result()
	return
}
