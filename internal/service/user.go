// @Author: 2014BDuck
// @Date: 2021/5/16

package service

import (
	"errors"
	"fmt"
	"github.com/2014bduck/entry-task/internal/constant"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserLoginRequest struct {
	Username string `form:"username" validate:"required,min=2,max=255"`
	Password string `form:"password" validate:"required,min=2,max=255"`
}

type UserRegisterRequest struct {
	Username   string `form:"username" validate:"required,min=2,max=255"`
	Password   string `form:"password" validate:"required,min=2,max=255"`
	Nickname   string `form:"nickname" validate:"required,min=2,max=255"`
	ProfilePic string `form:"profile_pic" validate:"-"` // Skip validation.
}

type UserLoginResponse struct {
	SessionID string
}

type UserRegisterResponse struct{}

func (svc *Service) UserLogin(param *UserLoginRequest) (*UserLoginResponse, error) {
	// Find user
	user, err := svc.dao.GetUserByName(param.Username)
	if err != nil {
		return nil, err
	}

	// Invalid cases
	if user.Password != param.Password {
		return nil, fmt.Errorf("svc.UserLogin: pwd incorrect")
	} else if user.Status != uint8(constant.EnabledStatus) {
		return nil, fmt.Errorf("svc.UserLogin: status disabled")
	}

	// Validation success
	// Setting session cache
	sessionID := uuid.NewV4()
	err = svc.cache.Cache.Set(sessionID.String(), []byte(param.Username))

	if err != nil {
		return nil, err
	}
	return &UserLoginResponse{SessionID: sessionID.String()}, nil
}

func (svc *Service) UserRegister(param *UserRegisterRequest) (*UserRegisterResponse, error) {
	// Validate username if existed
	_, err := svc.dao.GetUserByName(param.Username)
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("svc.UserRegister: username existed")
	}

	// Create User to DB
	_, err = svc.dao.CreateUser(param.Username, param.Password, param.Nickname, param.ProfilePic, uint8(constant.EnabledStatus))
	if err != nil {
		return nil, fmt.Errorf("svc.UserRegister: CreateUser error: %v", err)
	}

	return &UserRegisterResponse{}, nil
}

func (svc *Service) UserAuth(sessionID string) error {
	username, err := svc.cache.Cache.Get(sessionID)

	if err != nil || username == nil {
		return errors.New("svc.UserAuth failed")
	}
	return nil
}
