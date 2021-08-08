package erpc_service

import (
	"context"
	"errors"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	"github.com/2014bduck/entry-task/pkg/setting"
	"github.com/2014bduck/entry-task/pkg/upload"
	erpc_proto "github.com/2014bduck/entry-task/proto/erpc-proto"
	"github.com/agiledragon/gomonkey"
	"os"
	"reflect"
	"testing"
)

func TestUploadService_UploadFile(t *testing.T) {
	svc := NewUploadService(context.Background())
	erpcServer := erpc.NewServer(":8000")

	// Mock stuffs
	fileName := "test.png"
	outputFileName := "hashtest.png"
	serverUrl := "127.0.0.1"

	// Input
	request := erpc_proto.UploadRequest{
		FileType: uint32(upload.TypeImage),
		FileName: fileName,
		Content: make([]byte, 32),
	}

	t.Run("normal service register", func(t *testing.T){
		svc.RegisterUploadService(erpcServer)
	})

	t.Run("normal upload file", func(t *testing.T) {
		want := erpc_proto.UploadReply{
			FileUrl: serverUrl +"/"+ outputFileName,
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
		resp, err := svc.UploadFile(request)
		if err != nil {
			t.Errorf("TestUserService_UploadFile got error %v", err)
		}

		if reflect.DeepEqual(want, resp) {
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
		_, err := svc.UploadFile(request)
		if err == nil {
			t.Errorf("TestUserService_UploadFile should return error but didn't")
		}
	})
}