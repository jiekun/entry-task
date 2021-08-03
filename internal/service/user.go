// @Author: 2014BDuck
// @Date: 2021/5/16

package service

import (
	"errors"
	"github.com/2014bduck/entry-task/internal/constant"
	rpcproto "github.com/2014bduck/entry-task/internal/rpc-proto"
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
	var callLogin func(rpcproto.UserLoginRequest) (rpcproto.UserLoginResponse, error)
	svc.rpcClient.Call("Login", &callLogin)

	loginResp, err := callLogin(rpcproto.UserLoginRequest{
		Username: param.Username,
		Password: param.Password,
	})
	if err != nil {
		return nil, err
	}

	return &UserLoginResponse{loginResp.SessionID}, nil
}

func (svc *Service) CallRegister(param *UserRegisterRequest) (*UserRegisterResponse, error) {
	var callRegister func(rpcproto.UserRegisterRequest) (rpcproto.UserRegisterResponse, error)
	svc.rpcClient.Call("Register", &callRegister)

	_, err := callRegister(rpcproto.UserRegisterRequest{
		Username: param.Username,
		Password: param.Password,
		Nickname: param.Nickname,
	})
	if err != nil {
		return nil, err
	}

	return &UserRegisterResponse{}, nil
}

func (svc *Service) CallEditUser(param *UserEditRequest) (*UserEditResponse, error) {
	var callEdit func(rpcproto.UserEditRequest) (rpcproto.UserEditResponse, error)
	svc.rpcClient.Call("EditUser", &callEdit)

	_, err := callEdit(rpcproto.UserEditRequest{
		SessionID:  param.SessionID,
		Nickname:   param.Nickname,
		ProfilePic: param.ProfilePic,
	})
	if err != nil {
		return nil, err
	}

	return &UserEditResponse{}, nil
}

func (svc *Service) CallGetUser(param *UserGetRequest) (*UserGetResponse, error) {
	var callGet func(rpcproto.UserGetRequest) (rpcproto.UserGetResponse, error)
	svc.rpcClient.Call("GetUser", &callGet)

	getResp, err := callGet(rpcproto.UserGetRequest{
		SessionID: param.SessionID,
	})
	if err != nil {
		return nil, err
	}

	return &UserGetResponse{
		getResp.Username,
		getResp.Nickname,
		getResp.ProfilePic,
	}, nil
}

func (svc *Service) UserAuth(sessionID string) (string, error) {
	username, err := svc.cache.Cache.Get(constant.SessionIDCachePrefix + sessionID)

	if err != nil || username == nil {
		return "", errors.New("svc.UserAuth failed")
	}
	return string(username), nil
}
