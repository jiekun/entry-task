// @Author: 2014BDuck
// @Date: 2021/8/3

package erpc_service

import (
	"context"
	"errors"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/constant"
	"github.com/2014bduck/entry-task/internal/dao"
	"github.com/2014bduck/entry-task/pkg/upload"
	erpc_proto "github.com/2014bduck/entry-task/proto/erpc-proto"
	"os"
)

type UploadService struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.RedisCache
}

func NewUploadService(ctx context.Context) UploadService {
	svc := UploadService{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)

	return svc
}

func (svc UploadService) UploadFile(r *erpc_proto.UploadRequest) (*erpc_proto.UploadReply, error) {
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
	return &erpc_proto.UploadReply{FileUrl: fileUrl, FileName: fileName}, nil

}

func (svc UserService) UserAuth(sessionID string) (string, error) {
	username, err := svc.cache.Cache.Get(svc.ctx, constant.SessionIDCachePrefix+sessionID).Result()

	if err != nil {
		return "", errors.New("svc.UserAuth failed")
	}
	return username, nil
}
