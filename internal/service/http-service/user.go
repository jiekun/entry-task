// @Author: 2014BDuck
// @Date: 2021/5/16

package http_service

import (
	"errors"
	"github.com/jiekun/entry-task/internal/constant"
	"github.com/jiekun/entry-task/proto"
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

var userClient proto.UserServiceClient

func (svc *Service) Login(param *LoginRequest) (*LoginResponse, error) {
	resp, err := svc.getUserClient().Login(svc.ctx, &proto.LoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}
	return &LoginResponse{SessionID: resp.SessionId}, nil
}

func (svc *Service) Register(param *RegisterUserRequest) (*RegisterUserResponse, error) {
	_, err := svc.getUserClient().Register(svc.ctx, &proto.RegisterRequest{
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

func (svc *Service) EditUser(param *EditUserRequest) (*EditUserResponse, error) {
	_, err := svc.getUserClient().EditUser(svc.ctx, &proto.EditUserRequest{
		SessionId:  param.SessionID,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}
	return &EditUserResponse{}, nil
}

func (svc *Service) GetUser(param *GetUserRequest) (*GetUserResponse, error) {
	resp, err := svc.getUserClient().GetUser(svc.ctx, &proto.GetUserRequest{
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

func (svc *Service) AuthUser(sessionID string) (string, error) {
	username, err := svc.cache.Cache.Get(svc.ctx, constant.SessionIDCachePrefix+sessionID).Result()

	if err != nil {
		return "", errors.New("svc.UserAuth failed")
	}
	return username, nil
}

func (svc *Service) getUserClient() proto.UserServiceClient {
	if userClient == nil {
		userClient = proto.NewUserServiceClient(svc.client)
	}
	return userClient
}
