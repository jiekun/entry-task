// @Author: 2014BDuck
// @Date: 2021/8/2

package api

import (
	"github.com/jiekun/entry-task/global"
	errcode "github.com/jiekun/entry-task/internal/error"
	"github.com/jiekun/entry-task/internal/service/http-service"
	"github.com/jiekun/entry-task/pkg/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (Upload) Upload(c *gin.Context) {
	response := resp.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType, _ := strconv.Atoi(c.PostForm("type"))
	if err != nil {
		global.Logger.Errorf("app.Upload errs: %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
	}

	if fileHeader == nil || fileType < 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := http_service.New(c.Request.Context())
	uploadResp, err := svc.UploadFile(fileType, file, fileHeader)
	if err != nil {
		global.Logger.Errorf("app.Upload errs: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadPicFailed)
		return
	}
	response.ToResponse("Upload File Succeed.", uploadResp)
	return
}
