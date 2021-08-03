// @Author: 2014BDuck
// @Date: 2021/7/11

package service

import (
	"context"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/dao"
	"google.golang.org/grpc"
)

type Service struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.InProcessCache
	rpcClient *grpc.ClientConn
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)
	svc.rpcClient = global.GRPCClient

	return svc
}
