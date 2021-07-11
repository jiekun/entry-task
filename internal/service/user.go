// @Author: 2014BDuck
// @Date: 2021/5/16

package service

import (
	"errors"
	"github.com/satori/go.uuid"
)

type UserLoginRequest struct {
	Username string `form:"username" validate:"required,min=2,max=255"`
	Password string `form:"password" validate:"required,min=2,max=255"`
}

type UserLoginResponse struct {
	SessionID string
}

func (svc *Service) UserLogin(param *UserLoginRequest) (*UserLoginResponse, error) {
	// Validation
	err := svc.dao.ValidateUser(param.Username, param.Password)
	if err != nil {
		return nil, err
	}

	// Setting session cache
	sessionID := uuid.NewV4()
	err = svc.cache.Cache.Set(sessionID.String(), []byte(param.Username))

	if err != nil {
		return nil, err
	}
	return &UserLoginResponse{SessionID: sessionID.String()}, nil
}

func (svc *Service) UserAuth(sessionID string) error {
	username, err := svc.cache.Cache.Get(sessionID)

	if err != nil || username == nil {
		return errors.New("svc.UserAuth failed.")
	}
	return nil
}
