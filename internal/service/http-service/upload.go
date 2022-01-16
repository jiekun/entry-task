// @Author: 2014BDuck
// @Date: 2021/8/2

package http_service

import (
	"errors"
	"github.com/2014bduck/entry-task/pkg/upload"
	"github.com/2014bduck/entry-task/proto"
	"mime/multipart"
)

type UploadFileRequest struct {
}

type UploadFileResponse struct {
	FileName string `json:"file_name"`
	FileUrl  string `json:"file_url"`
}

var uploadClient proto.UploadServiceClient

func (svc *Service) UploadFile(fileType int, file multipart.File, fileHeader *multipart.FileHeader) (*UploadFileResponse, error) {
	// Basic Check
	if !upload.CheckContainExt(upload.FileType(fileType), fileHeader.Filename) {
		return nil, errors.New("svc.UploadFile: file suffix is not supported")
	}

	if upload.CheckMaxSize(upload.FileType(fileType), file) {
		return nil, errors.New("svc.UploadFile: exceed maximum file limit")
	}

	// Read to []byte
	content, err := upload.GetFileByte(fileHeader)
	if err != nil {
		return nil, errors.New("svc.UploadFile: failed reading file to []byte")
	}

	// Transfer []byte via RPC
	resp, err := svc.getUploadClient().UploadFile(svc.ctx, &proto.UploadRequest{
		FileType: uint32(fileType),
		FileName: fileHeader.Filename,
		Content:  content,
	})
	if err != nil {
		return nil, err
	}
	return &UploadFileResponse{
		FileName: resp.FileName,
		FileUrl:  resp.FileUrl,
	}, nil

}

func (svc *Service) getUploadClient() proto.UploadServiceClient {
	if uploadClient == nil {
		uploadClient = proto.NewUploadServiceClient(svc.client)
	}
	return uploadClient
}
