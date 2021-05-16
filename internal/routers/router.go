// @Author: 2014BDuck
// @Date: 2021/5/16

package routers

import (
	"github.com/2014bduck/entry-task/internal/routers/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	ping := api.NewPing()
	apiGroup := r.Group("/api/")
	{
		apiGroup.GET("/ping", ping.Ping)
	}
	return r
}
