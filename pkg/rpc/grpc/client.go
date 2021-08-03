// @Author: 2014BDuck
// @Date: 2021/8/3

package grpc

import (
	"context"
	"google.golang.org/grpc"
)

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	return grpc.DialContext(ctx, target, opts...)
}
