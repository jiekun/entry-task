// @Author: 2014BDuck
// @Date: 2021/5/16

package api

import (
	"fmt"
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
	param := service.LoginRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Errorf("app.Login errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	loginResponse, err := svc.CallLogin(&param)
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
	param := service.RegisterUserRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	loginResponse, err := svc.CallRegister(&param)
	if err != nil {
		global.Logger.Errorf("app.Register errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserRegister)
		return
	}
	response.ToResponse("Register Succeed.", loginResponse)
	return
}

func (u User) Edit(c *gin.Context) {
	response := resp.NewResponse(c)
	param := service.EditUserRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Errorf("app.Edit errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	// Get sessionID set by Auth middleware
	sessionID, _ := c.Get("sessionID")
	param.SessionID = fmt.Sprintf("%v", sessionID)

	svc := service.New(c.Request.Context())
	editResponse, err := svc.CallEditUser(&param)
	if err != nil {
		global.Logger.Errorf("app.Edit errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserRegister)
		return
	}
	response.ToResponse("Edit Succeed.", editResponse)
	return
}

func (u User) Get(c *gin.Context) {
	response := resp.NewResponse(c)
	param := service.GetUserRequest{}

	// Get sessionID set by Auth middleware
	sessionID, _ := c.Get("sessionID")
	param.SessionID = fmt.Sprintf("%v", sessionID)

	svc := service.New(c.Request.Context())
	getResponse, err := svc.CallGetUser(&param)
	if err != nil {
		global.Logger.Errorf("app.Get errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserGet)
		return
	}
	response.ToResponse("Get Succeed.", getResponse)
	return
}