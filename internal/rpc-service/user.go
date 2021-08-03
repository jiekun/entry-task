// @Author: 2014BDuck
// @Date: 2021/5/16

package rpcservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/constant"
	rpcproto "github.com/2014bduck/entry-task/internal/rpc-proto"
	"github.com/2014bduck/entry-task/pkg/hashing"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)


func (svc *Service) UserLogin(param rpcproto.UserLoginRequest) (*rpcproto.UserLoginResponse, error) {
	// Find user
	user, err := svc.dao.GetUserByName(param.Username)
	if err != nil {
		return &rpcproto.UserLoginResponse{}, err
	}

	// Invalid cases
	hashedPass := hashing.HashPassword(param.Password)
	if user.Password != hashedPass {
		return &rpcproto.UserLoginResponse{}, errors.New("svc.UserLogin: pwd incorrect")
	} else if user.Status != uint8(constant.EnabledStatus) {
		return &rpcproto.UserLoginResponse{}, errors.New("svc.UserLogin: status disabled")
	}

	// Validation success
	// Setting session cache
	sessionID := uuid.NewV4()
	err = svc.cache.Cache.Set(constant.SessionIDCachePrefix+sessionID.String(), []byte(param.Username))

	if err != nil {
		return &rpcproto.UserLoginResponse{}, err
	}
	return &rpcproto.UserLoginResponse{SessionID: sessionID.String()}, nil
}

func (svc *Service) UserRegister(param rpcproto.UserRegisterRequest) (*rpcproto.UserRegisterResponse, error) {
	// Validate username if existed
	_, err := svc.dao.GetUserByName(param.Username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("svc.UserRegister: username existed")
	}

	// Add Salt to pass
	hashedPass := hashing.HashPassword(param.Password)

	// Create User to DB
	_, err = svc.dao.CreateUser(param.Username, hashedPass, param.Nickname, param.ProfilePic, uint8(constant.EnabledStatus))
	if err != nil {
		return &rpcproto.UserRegisterResponse{}, fmt.Errorf("svc.UserRegister: CreateUser error: %v", err)
	}

	return &rpcproto.UserRegisterResponse{}, nil
}

func (svc *Service) UserEdit(param rpcproto.UserEditRequest) (*rpcproto.UserEditResponse, error) {
	// Get Username
	username, err := svc.UserAuth(param.SessionID)
	if err != nil {
		return &rpcproto.UserEditResponse{}, err
	}

	// Query current user
	user, err := svc.dao.GetUserByName(username)
	if err != nil {
		return &rpcproto.UserEditResponse{}, fmt.Errorf("svc.UserEdit: %v", err)
	}

	// Validate user status
	if constant.Status(user.Status) != constant.EnabledStatus {
		return nil, errors.New("svc.UserEdit: Invalid user status")
	}

	// Update user data
	err = svc.dao.UpdateUser(user.ID, param.Nickname, param.ProfilePic)
	if err != nil {
		return nil, fmt.Errorf("svc.UserEdit: %v", err)
	}
	return &rpcproto.UserEditResponse{}, nil
}

func (svc *Service) UserGet(param rpcproto.UserGetRequest) (*rpcproto.UserGetResponse, error) {
	// Get Username
	username, err := svc.UserAuth(param.SessionID)
	if err != nil {
		return &rpcproto.UserGetResponse{}, err
	}

	cacheKey := constant.UserProfileCachePrefix + username

	// Try loading user info from cache
	userProfCache, err := svc.cache.Cache.Get(cacheKey)
	if err == nil && userProfCache != nil {
		userGetCacheResp := rpcproto.UserGetResponse{}
		err = json.Unmarshal(userProfCache, &userGetCacheResp)
		if err != nil {
			global.Logger.Errorf("svc.UserGet: Unmarshal cache failed: %v", err)
		} else {
			return &userGetCacheResp, nil
		}
	}

	// Query user from DB
	user, err := svc.dao.GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("svc.UserGet: %v", err)
	}
	userGetResp := &rpcproto.UserGetResponse{
		Username:   user.Name,
		Nickname:   user.Nickname,
		ProfilePic: user.ProfilePic,
	}

	// Set user to cache
	cacheUser, _ := json.Marshal(userGetResp)
	err = svc.cache.Cache.Set(cacheKey, cacheUser) // Omit error
	if err != nil {
		global.Logger.Errorf("svc.UserGet: set cache failed: %v", err)
	}

	return userGetResp, nil
}

func (svc *Service) UserAuth(sessionID string) (string, error) {
	username, err := svc.cache.Cache.Get(constant.SessionIDCachePrefix + sessionID)

	if err != nil || username == nil {
		return "", errors.New("svc.UserAuth failed")
	}
	return string(username), nil
}
