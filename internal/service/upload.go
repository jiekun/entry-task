// @Author: 2014BDuck
// @Date: 2021/8/2

package service

import (
	"errors"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/pkg/upload"
	"mime/multipart"
	"os"
)

type UploadFileRequest struct {
}

type UploadFileResponse struct {
	FileName string `json:"file_name"`
	FileUrl  string `json:"file_url"`
}

func (svc *Service) UploadFile(fileType int, file multipart.File, fileHeader *multipart.FileHeader) (*UploadFileResponse, error) {
	fileName := upload.GetFileName(fileHeader.Filename) // MD5'd
	uploadSavePath := upload.GetSavePath()
	dst := uploadSavePath + "/" + fileName
	if !upload.CheckContainExt(upload.FileType(fileType), fileName) {
		return nil, errors.New("svc.UploadFile: file suffix is not supported")
	}
	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(dst, os.ModePerm); err != nil {
			return nil, errors.New("svc.UploadFile: failed to create save directory")
		}
	}
	if upload.CheckMaxSize(upload.FileType(fileType), file) {
		return nil, errors.New("svc.UploadFile: exceed maximum file limit")
	}
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("svc.UploadFile: insufficient file permissions")
	}
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}
	fileUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &UploadFileResponse{FileUrl: fileUrl, FileName: fileName}, nil
}
