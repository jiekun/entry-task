// @Author: 2014BDuck
// @Date: 2021/7/11

package dao

import (
	"github.com/allegro/bigcache/v3"
	"time"
)

type InProcessCache struct {
	Cache *bigcache.BigCache
}

func NewCache() *InProcessCache {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(30 * time.Minute))
	return &InProcessCache{Cache: cache}
}
