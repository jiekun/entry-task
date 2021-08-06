// @Author: 2014BDuck
// @Date: 2021/5/16

package service

import (
	"encoding/gob"
	"errors"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/constant"
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	erpc_proto "github.com/2014bduck/entry-task/proto/erpc-proto"
	pb "github.com/2014bduck/entry-task/proto/grpc-proto"
)

type LoginRequest struct {
	Username string `form:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" binding:"required,min=6,max=64"`
}

type RegisterUserRequest struct {
	Username   string `form:"username" binding:"required,min=6,max=64"`
	Password   string `form:"password" binding:"required,min=6,max=64"`
	Nickname   string `form:"nickname" binding:"required,min=6,max=64"`
	ProfilePic string `form:"profile_pic" binding:"-"` // Skip validation.
}

type EditUserRequest struct {
	SessionID  string `form:"session_id"`
	Nickname   string `form:"nickname" binding:"min=6,max=64"`
	ProfilePic string `form:"profile_pic" binding:"max=1024"`
}

type GetUserRequest struct {
	SessionID string `form:"session_id"`
}

type LoginResponse struct {
	SessionID string `json:"session_id"`
}

type RegisterUserResponse struct{}

type EditUserResponse struct{}

type GetUserResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"`
}

func RegisterUserServiceProto() {
	gob.Register(erpc_proto.LoginRequest{})
	gob.Register(erpc_proto.LoginReply{})
	gob.Register(erpc_proto.RegisterRequest{})
	gob.Register(erpc_proto.RegisterReply{})
	gob.Register(erpc_proto.EditUserRequest{})
	gob.Register(erpc_proto.EditUserReply{})
	gob.Register(erpc_proto.GetUserRequest{})
	gob.Register(erpc_proto.GetUserReply{})
}

func (svc *Service) CallLogin(param *LoginRequest) (*LoginResponse, error) {
	RPCLogin := erpc_proto.Login
	c, err := svc.getConn()
	if err != nil{
		return nil, err
	}
	c.Call("Login", &RPCLogin)

	resp, err := RPCLogin(erpc_proto.LoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LoginResponse{SessionID: resp.SessionId}, nil
}

func (svc *Service) CallRegister(param *RegisterUserRequest) (*RegisterUserResponse, error) {
	RPCRegister := erpc_proto.Register
	c, err := svc.getConn()
	if err != nil{
		return nil, err
	}
	c.Call("Register", &RPCRegister)
	_, err = RPCRegister(erpc_proto.RegisterRequest{
		Username:   param.Username,
		Password:   param.Password,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &RegisterUserResponse{}, nil
}

func (svc *Service) CallEditUser(param *EditUserRequest) (*EditUserResponse, error) {
	RPCEditUser := erpc_proto.EditUser
	c, err := svc.getConn()
	if err != nil{
		return nil, err
	}
	c.Call("EditUser", &RPCEditUser)
	_, err = RPCEditUser(erpc_proto.EditUserRequest{
		SessionId:  param.SessionID,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &EditUserResponse{}, nil
}

func (svc *Service) CallGetUser(param *GetUserRequest) (*GetUserResponse, error) {
	RPCGetUser := erpc_proto.GetUser
	c, err := svc.getConn()
	if err != nil{
		return nil, err
	}
	c.Call("GetUser", &RPCGetUser)
	resp, err := RPCGetUser(erpc_proto.GetUserRequest{
		SessionId: param.SessionID,
	})
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{
		Username:   resp.Username,
		Nickname:   resp.Nickname,
		ProfilePic: resp.ProfilePic,
	}, nil
}

func (svc *Service) CallLoginGRpc(param *LoginRequest) (*LoginResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.gRpcClient)
	resp, err := userServiceClient.Login(svc.ctx, &pb.LoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LoginResponse{SessionID: resp.SessionId}, nil
}

func (svc *Service) CallRegisteGRpc(param *RegisterUserRequest) (*RegisterUserResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.gRpcClient)
	_, err := userServiceClient.Register(svc.ctx, &pb.RegisterRequest{
		Username:   param.Username,
		Password:   param.Password,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &RegisterUserResponse{}, nil
}

func (svc *Service) CallEditUserGRpc(param *EditUserRequest) (*EditUserResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.gRpcClient)
	_, err := userServiceClient.EditUser(svc.ctx, &pb.EditUserRequest{
		SessionId:  param.SessionID,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &EditUserResponse{}, nil
}

func (svc *Service) CallGetUserGRpc(param *GetUserRequest) (*GetUserResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.gRpcClient)
	resp, err := userServiceClient.GetUser(svc.ctx, &pb.GetUserRequest{
		SessionId: param.SessionID,
	})
	if err != nil {
		return nil, err
	}
	return &GetUserResponse{
		Username:   resp.Username,
		Nickname:   resp.Nickname,
		ProfilePic: resp.ProfilePic,
	}, nil
}

func (svc *Service) UserAuth(sessionID string) (string, error) {
	username, err := svc.cache.Cache.Get(svc.ctx, constant.SessionIDCachePrefix+sessionID).Result()

	if err != nil {
		return "", errors.New("svc.UserAuth failed")
	}
	return username, nil
}

func (svc *Service) getConn() (*erpc.Client, error) {
	if global.RPCClientPool == nil {
		cp, err := erpc.NewConnectionPool(global.ClientSetting.RPCHost, global.ClientSetting.ConnNum)
		if err != nil {
			return nil, err
		}
		global.RPCClientPool = cp
	}
	conn, lock, err := global.RPCClientPool.Get()
	if err != nil {
		return nil, err
	}
	return &erpc.Client{Conn: *conn, Lock: lock}, nil
}
