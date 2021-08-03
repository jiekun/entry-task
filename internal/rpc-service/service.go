// @Author: 2014BDuck
// @Date: 2021/7/11

package rpcservice

import (
	"context"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/dao"
)

type Service struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.InProcessCache
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)

	return svc
}
