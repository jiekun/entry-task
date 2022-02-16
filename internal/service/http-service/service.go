// @Author: 2014BDuck
// @Date: 2021/7/11

package http_service

import (
	"context"
	"github.com/jiekun/entry-task/global"
	"github.com/jiekun/entry-task/internal/dao"
	"google.golang.org/grpc"
)

type Service struct {
	ctx    context.Context
	dao    *dao.Dao
	cache  *dao.RedisCache
	client *grpc.ClientConn
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)
	svc.client = global.GRPCClient

	return svc
}
