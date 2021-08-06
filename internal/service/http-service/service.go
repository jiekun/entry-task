// @Author: 2014BDuck
// @Date: 2021/7/11

package http_service

import (
	"context"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/dao"
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	"google.golang.org/grpc"
)

type Service struct {
	ctx          context.Context
	dao          *dao.Dao
	cache        *dao.RedisCache
	gRpcClient   *grpc.ClientConn
	eRpcConnPool *erpc.ConnectionPool
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)
	svc.gRpcClient = global.GRPCClient
	svc.eRpcConnPool = global.RPCClientPool

	return svc
}
