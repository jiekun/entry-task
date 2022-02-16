// @Author: 2014BDuck
// @Date: 2021/7/11

package middleware

import (
	errcode "github.com/jiekun/entry-task/internal/error"
	"github.com/jiekun/entry-task/internal/service/http-service"
	"github.com/jiekun/entry-task/pkg/resp"
	"github.com/gin-gonic/gin"
)

func SessionRequired(c *gin.Context) {
	response := resp.NewResponse(c)
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		// Abort the request with the appropriate error code
		response.ToAbortErrorResponse(errcode.ErrorUserNotLogin)
		return
	}else{
		// Continue down the chain to handler etc
		c.Set("sessionID", sessionID)
		c.Next()
	}
}

func LoginRequired(c *gin.Context) {
	response := resp.NewResponse(c)
	sessionID, err := c.Cookie("session_id")
	svc := http_service.New(c.Request.Context())
	username, err := svc.AuthUser(sessionID)
	if err != nil {
		// Abort the request with the appropriate error code
		response.ToAbortErrorResponse(errcode.ErrorUserNotLogin)
		return
	}else{
		// Continue down the chain to handler etc
		c.Set("username", username)
		c.Next()
	}
}

