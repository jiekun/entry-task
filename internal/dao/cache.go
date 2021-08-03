// @Author: 2014BDuck
// @Date: 2021/7/11

package dao

import (
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	Cache *redis.Client
}


func NewCache(cacheClient *redis.Client) *RedisCache {
	return &RedisCache{Cache: cacheClient}
}