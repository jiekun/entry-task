package grpc_service

import (
	"context"
	"errors"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/pkg/setting"
	"github.com/2014bduck/entry-task/pkg/upload"
	"github.com/2014bduck/entry-task/proto"
	"github.com/agiledragon/gomonkey"
	"os"
	"testing"
)

func TestUploadService_UploadFile(t *testing.T) {
	svc := NewUploadService(context.Background())

	// Mock stuffs
	fileName := "test.png"
	outputFileName := "hashtest.png"
	serverUrl := "127.0.0.1"

	// Input
	request := &proto.UploadRequest{
		FileType: uint32(upload.TypeImage),
		FileName: fileName,
		Content:  make([]byte, 32),
	}

	t.Run("normal upload file", func(t *testing.T) {
		want := &proto.UploadReply{
			FileUrl:  serverUrl + "/" + outputFileName,
			FileName: outputFileName,
		}
		// Mock DAO call
		patches := gomonkey.ApplyFunc(upload.GetFileName, func(string) string {
			return outputFileName
		})
		defer patches.Reset()
		patches.ApplyFunc(upload.GetSavePath, func() string {
			return ""
		})
		patches.ApplyFunc(upload.CheckSavePath, func(string) bool {
			return false
		})
		patches.ApplyFunc(upload.CheckPermission, func(string) bool {
			return false
		})
		patches.ApplyFunc(upload.SaveFileByte, func(*[]byte, string) error {
			return nil
		})
		patches.ApplyGlobalVar(&global.AppSetting, &setting.AppSettingS{
			UploadServerUrl: serverUrl,
		})
		// Test and compare
		resp, err := svc.UploadFile(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_UploadFile got error %v", err)
		}

		if want.FileName != resp.GetFileName() || resp.FileUrl != resp.GetFileUrl() {
			t.Errorf("TestUserService_UploadFile want %v got %v", want, resp)
		}
	})
	t.Run("upload file failed", func(t *testing.T) {

		// Mock DAO call
		patches := gomonkey.ApplyFunc(upload.GetFileName, func(string) string {
			return outputFileName
		})
		defer patches.Reset()
		patches.ApplyFunc(upload.GetSavePath, func() string {
			return ""
		})
		patches.ApplyFunc(upload.CheckSavePath, func(string) bool {
			return true
		})
		patches.ApplyFunc(upload.CreateSavePath, func(string, os.FileMode) error {
			return errors.New("errors")
		})
		// Test and compare
		_, err := svc.UploadFile(context.Background(), request)
		if err == nil {
			t.Errorf("TestUserService_UploadFile should return error but didn't")
		}
	})
}
