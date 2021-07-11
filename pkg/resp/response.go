// @Author: 2014BDuck
// @Date: 2021/7/11

package resp

import (
	errcode "github.com/2014bduck/entry-task/internal/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(msg string, data interface{}) {
	data = gin.H{"retcode": 0, "msg": "success", "data": data}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"retcode": err.GetCode(), "msg": err.GetMsg()}
	details := err.GetDetails()
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}

func (r *Response) ToAbortErrorResponse(err *errcode.Error) {
	response := gin.H{"retcode": err.GetCode(), "msg": err.GetMsg()}
	details := err.GetDetails()
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.AbortWithStatusJSON(err.StatusCode(), response)
}
