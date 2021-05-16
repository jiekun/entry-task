// @Author: 2014BDuck
// @Date: 2021/5/16

package api

import "github.com/gin-gonic/gin"

type User struct{}

func NewUser() User {
	return User{}
}

func (u User) Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}