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
		return
	}
	svc := service.New(c.Request.Context())
	loginResponse, err := svc.UserLogin(&param)
	if err != nil {
		global.Logger.Errorf("app.UserLogin errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUserLogin)
		return

	}
	response.ToResponse("Login Succeed.", loginResponse)
	return
}
