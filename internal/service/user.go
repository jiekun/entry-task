// @Author: 2014BDuck
// @Date: 2021/5/16

package service

import (
	"errors"
	"github.com/2014bduck/entry-task/internal/constant"
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

func (svc *Service) CallLogin(param *LoginRequest) (*LoginResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
	resp, err := userServiceClient.Login(svc.ctx, &pb.LoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LoginResponse{SessionID: resp.SessionId}, nil
}

func (svc *Service) CallRegister(param *RegisterUserRequest) (*RegisterUserResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
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

func (svc *Service) CallEditUser(param *EditUserRequest) (*EditUserResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
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

func (svc *Service) CallGetUser(param *GetUserRequest) (*GetUserResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
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
