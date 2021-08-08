// @Author: 2014BDuck
// @Date: 2021/7/11

package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	Cache *redis.Client
}

func NewCache(cacheClient *redis.Client) *RedisCache {
	return &RedisCache{Cache: cacheClient}
}

func (cache *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := cache.Cache.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (cache *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := cache.Cache.Set(ctx, key, value, expiration).Err()
	return err
}
