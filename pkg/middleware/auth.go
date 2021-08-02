// @Author: 2014BDuck
// @Date: 2021/7/11

package middleware

import (
	errcode "github.com/2014bduck/entry-task/internal/error"
	"github.com/2014bduck/entry-task/internal/service"
	"github.com/2014bduck/entry-task/pkg/resp"
	"github.com/gin-gonic/gin"
)

func SessionRequired(c *gin.Context) {
	response := resp.NewResponse(c)
	sessionID, err := c.Cookie("session_id")
	svc := service.New(c.Request.Context())
	err = svc.UserAuth(sessionID)
	if err != nil {
		// Abort the request with the appropriate error code
		response.ToAbortErrorResponse(errcode.ErrorUserNotLogin)
		return
	}else{
		// Continue down the chain to handler etc
		c.Next()
	}
}