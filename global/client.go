// @Author: 2014BDuck
// @Date: 2021/8/3

package global

import (
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	"google.golang.org/grpc"
)

var (
	GRPCClient    *grpc.ClientConn
	RPCClientPool *erpc.ConnectionPool
)
