// @Author: 2014BDuck
// @Date: 2021/5/16

package routers

import (
	"github.com/2014bduck/entry-task/internal/routers/api"
	"github.com/2014bduck/entry-task/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	ping := api.NewPing()
	user := api.NewUser()
	apiGroup := r.Group("/api/")
	{
		apiGroup.POST("/user/login", user.Login)
		apiGroup.POST("/user/register", user.Redgister)
	}

	authGroup := r.Group("/api/")
	authGroup.Use(middleware.SessionRequired)
	{
		authGroup.GET("/ping", ping.Ping)
	}
	return r
}
