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
	upload := api.NewUpload()

	// No need login
	apiGroup := r.Group("/api/")
	{
		apiGroup.POST("/user/login", user.Login)
		apiGroup.POST("/user/register", user.Register)
		apiGroup.GET("/ping", ping.Ping)
	}

	// Session Required
	sessionGroup := r.Group("/api/")
	sessionGroup.Use(middleware.SessionRequired)
	{
		sessionGroup.GET("/user/get", user.Get)
		sessionGroup.POST("/user/edit", user.Edit)
	}

	// Login Required
	loginGroup := r.Group("/api/")
	loginGroup.Use(middleware.LoginRequired)
	{
		apiGroup.POST("/upload/file", upload.Upload)
	}

	return r
}
