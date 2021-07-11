// @Author: 2014BDuck
// @Date: 2021/7/11

package dao

import (
	"github.com/allegro/bigcache/v3"
)

type InProcessCache struct {
	Cache *bigcache.BigCache
}


func NewCache(cacheClient *bigcache.BigCache) *InProcessCache {
	return &InProcessCache{Cache: cacheClient}
}