// @Author: 2014BDuck
// @Date: 2021/5/16

package service

import (
	"errors"
	"github.com/2014bduck/entry-task/internal/constant"
	pb "github.com/2014bduck/entry-task/proto/grpc-proto"
)

type UserLoginRequest struct {
	Username string `form:"username" binding:"required,min=2,max=255"`
	Password string `form:"password" binding:"required,min=2,max=255"`
}

type UserRegisterRequest struct {
	Username   string `form:"username" binding:"required,min=2,max=255"`
	Password   string `form:"password" binding:"required,min=2,max=255"`
	Nickname   string `form:"nickname" binding:"required,min=2,max=255"`
	ProfilePic string `form:"profile_pic" binding:"-"` // Skip validation.
}

type UserEditRequest struct {
	SessionID  string `form:"session_id"`
	Nickname   string `form:"nickname"`
	ProfilePic string `form:"profile_pic"`
}

type UserGetRequest struct {
	SessionID string `form:"session_id"`
}

type UserLoginResponse struct {
	SessionID string `json:"session_id"`
}

type UserRegisterResponse struct{}

type UserEditResponse struct{}

type UserGetResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"`
}

func (svc *Service) CallLogin(param *UserLoginRequest) (*UserLoginResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
	resp, err := userServiceClient.Login(svc.ctx, &pb.LoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}
	return &UserLoginResponse{SessionID: resp.SessionId}, nil
}

func (svc *Service) CallRegister(param *UserRegisterRequest) (*UserRegisterResponse, error) {
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
	return &UserRegisterResponse{}, nil
}

func (svc *Service) CallEditUser(param *UserEditRequest) (*UserEditResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
	_, err := userServiceClient.EditUser(svc.ctx, &pb.EditUserRequest{
		SessionId:  param.SessionID,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &UserEditResponse{}, nil
}

func (svc *Service) CallGetUser(param *UserGetRequest) (*UserGetResponse, error) {
	userServiceClient := pb.NewUserServiceClient(svc.rpcClient)
	resp, err := userServiceClient.GetUser(svc.ctx, &pb.GetUserRequest{
		SessionId: param.SessionID,
	})
	if err != nil {
		return nil, err
	}
	return &UserGetResponse{
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
