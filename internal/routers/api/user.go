// @Author: 2014BDuck
// @Date: 2021/5/16

package api

import (
	"github.com/2014bduck/entry-task/global"
	errcode "github.com/2014bduck/entry-task/internal/error"
	"github.com/2014bduck/entry-task/internal/service"
	"github.com/2014bduck/entry-task/pkg/resp"
	"github.com/gin-gonic/gin"
)

type User struct{}

func NewUser() User {
	return User{}
}

func (u User) Login(c *gin.Context) {
	response := resp.NewResponse(c)
	param := service.UserLoginRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Errorf("app.Login errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	loginResponse, err := svc.UserLogin(&param)
	if err != nil {
		global.Logger.Errorf("app.Login errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserLogin)
		return
	}
	c.SetCookie("session_id", loginResponse.SessionID, 3600, "/", "", false, true)
	response.ToResponse("Login Succeed.", loginResponse)
	return
}

func (u User) Register(c *gin.Context) {
	response := resp.NewResponse(c)
	param := service.UserRegisterRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	loginResponse, err := svc.UserRegister(&param)
	if err != nil {
		global.Logger.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserRegister)
		return
	}
	response.ToResponse("Register Succeed.", loginResponse)
	return
}
