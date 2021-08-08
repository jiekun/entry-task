// @Author: 2014BDuck
// @Date: 2021/8/6

package erpc_service

import (
	"context"
	"errors"
	"github.com/2014bduck/entry-task/internal/dao"
	"github.com/2014bduck/entry-task/internal/models"
	"github.com/2014bduck/entry-task/pkg/hashing"
	erpc_proto "github.com/2014bduck/entry-task/proto/erpc-proto"
	"github.com/agiledragon/gomonkey"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestUserService_Register(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	nickname := "test_nickname"
	password := "test_password"

	// Input
	request := erpc_proto.RegisterRequest{
		Username: username,
		Nickname: nickname,
		Password: password,
	}

	t.Run("normal register", func(t *testing.T) {
		// Target output
		want := &erpc_proto.RegisterReply{}

		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{}, gorm.ErrRecordNotFound
		})
		defer patches.Reset()
		patches.ApplyFunc(hashing.HashPassword, func(_ string) string {
			return "mock_hashing"
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "CreateUser", func(_ *dao.Dao, _, _, _, _ string, _ uint8) (*models.UserTab, error) {
			return &models.UserTab{
				Name:     username,
				Nickname: nickname,
				Password: password, // It's hashed actually.
			}, nil
		})

		// Test and compare with reflect.DeepEqual
		resp, err := svc.Register(request)
		if err != nil {
			t.Errorf("TestUserService_Register got error %v", err)
		}

		if !reflect.DeepEqual(want, resp) {
			t.Errorf("TestUserService_Register want: %v got %v", want, resp)
		}
	})

	t.Run("invalid register", func(t *testing.T) {
		// Mock GetUser with record found
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{}, nil
		})
		defer patches.Reset()

		// should return an err
		_, err := svc.Register(request)
		if err == nil {
			t.Error("TestUserService_Register should return error but didn't")
		}
	})

}

func TestUserService_Login(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	nickname := "test_nickname"
	password := "test_password"

	// Input
	request := erpc_proto.LoginRequest{
		Username: username,
		Password: password,
	}

	t.Run("normal login", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				Name:     username,
				Nickname: nickname,
				Password: password,
			}, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(hashing.HashPassword, func(_ string) string {
			return password
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface {}, _ time.Duration) error {
			return nil
		})

		// Test and compare
		resp, err := svc.Login(request)
		if err != nil {
			t.Errorf("TestUserService_Login got error %v", err)
		}

		if resp.SessionId == "" {
			t.Errorf("TestUserService_Login got %v", resp)
		}
	})

	t.Run("login no such user", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{}, gorm.ErrRecordNotFound
		})
		defer patches.Reset()
		// Test and compare
		_, err := svc.Login(request)
		if err == nil {
			t.Errorf("TestUserService_Login should return err but didn't")
		}
	})

	t.Run("login incorrect password", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				Name:     username,
				Nickname: nickname,
				Password: "",
			}, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(hashing.HashPassword, func(_ string) string {
			return password
		})
		// Test and compare
		_, err := svc.Login(request)
		if err == nil {
			t.Errorf("TestUserService_Login should return err but didn't")
		}
	})

	t.Run("login failed to set session", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				Name:     username,
				Nickname: nickname,
				Password: "",
			}, nil
		})
		defer patches.Reset()
		patches.ApplyFunc(hashing.HashPassword, func(_ string) string {
			return password
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface {}, _ time.Duration) error {
			return errors.New("failed")
		})

		// Test and compare
		_, err := svc.Login(request)
		if err == nil {
			t.Errorf("TestUserService_Login should return err but didn't")
		}
	})
}
