// @Author: 2014BDuck
// @Date: 2021/8/2

package http_service

import (
	"encoding/gob"
	"errors"
	"github.com/2014bduck/entry-task/pkg/upload"
	erpc_proto "github.com/2014bduck/entry-task/proto/erpc-proto"
	pb "github.com/2014bduck/entry-task/proto/grpc-proto"
	"mime/multipart"
)

type UploadFileRequest struct {
}

type UploadFileResponse struct {
	FileName string `json:"file_name"`
	FileUrl  string `json:"file_url"`
}

func RegisterUploadServiceProto() {
	gob.Register(erpc_proto.UploadRequest{})
	gob.Register(erpc_proto.UploadReply{})
}

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
	RPCUploadFile := erpc_proto.UploadFile
	c, err := svc.getConn()
	if err != nil{
		return nil, err
	}
	c.Call("UploadFile", &RPCUploadFile)
	resp, err := RPCUploadFile(erpc_proto.UploadRequest{
		FileType: uint32(fileType),
		FileName: fileHeader.Filename,
		Content: content,

	})
	if err != nil {
		return nil, err
	}
	return &UploadFileResponse{
		FileName: resp.FileName,
		FileUrl:  resp.FileUrl,
	}, nil
}


func (svc *Service) UploadFileGRpc(fileType int, file multipart.File, fileHeader *multipart.FileHeader) (*UploadFileResponse, error) {
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
	userServiceClient := pb.NewUploadServiceClient(svc.gRpcClient)
	resp, err := userServiceClient.UploadFile(svc.ctx, &pb.UploadRequest{
		FileType: uint32(fileType),
		FileName:   fileHeader.Filename,
		Content:    content,
	})
	if err != nil {
		return nil, err
	}
	return &UploadFileResponse{
		FileName: resp.FileName,
		FileUrl:  resp.FileUrl,
	}, nil

}
