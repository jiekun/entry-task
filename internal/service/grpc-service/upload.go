// @Author: 2014BDuck
// @Date: 2021/8/3

package grpc_service

import (
	"context"
	"errors"
	"github.com/jiekun/entry-task/global"
	"github.com/jiekun/entry-task/internal/dao"
	"github.com/jiekun/entry-task/pkg/upload"
	"github.com/jiekun/entry-task/proto"
	"os"
)

type UploadService struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.RedisCache
	proto.UnimplementedUploadServiceServer
}

func NewUploadService(ctx context.Context) UploadService {
	svc := UploadService{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)

	return svc
}

func (svc UploadService) UploadFile(ctx context.Context, r *proto.UploadRequest) (*proto.UploadReply, error) {
	fileName := upload.GetFileName(r.FileName) // MD5'd
	uploadSavePath := upload.GetSavePath()
	dst := uploadSavePath + "/" + fileName

	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(dst, os.ModePerm); err != nil {
			return nil, errors.New("svc.UploadFile: failed to create save directory")
		}
	}

	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("svc.UploadFile: insufficient file permissions")
	}
	if err := upload.SaveFileByte(&r.Content, dst); err != nil {
		return nil, err
	}
	fileUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &proto.UploadReply{FileUrl: fileUrl, FileName: fileName}, nil

}
