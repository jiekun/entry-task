// @Author: 2014BDuck
// @Date: 2021/8/6

package grpc_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jiekun/entry-task/internal/constant"
	"github.com/jiekun/entry-task/internal/dao"
	"github.com/jiekun/entry-task/internal/models"
	"github.com/jiekun/entry-task/pkg/hashing"
	"github.com/jiekun/entry-task/proto"
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
	request := &proto.RegisterRequest{
		Username: username,
		Nickname: nickname,
		Password: password,
	}

	t.Run("normal register", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{}, gorm.ErrRecordNotFound
		})
		defer patches.Reset()
		patches.ApplyFunc(hashing.HashPasswordBcrypt, func(_ string) (string, error) {
			return "mock_hashing", nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "CreateUser", func(_ *dao.Dao, _, _, _, _ string, _ uint8) (*models.UserTab, error) {
			return &models.UserTab{
				Name:     username,
				Nickname: nickname,
				Password: password, // It's hashed actually.
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface{}, _ time.Duration) error {
			return nil
		})

		// Test and compare with reflect.DeepEqual
		_, err := svc.Register(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Register got error %v", err)
		}
	})

	t.Run("invalid register", func(t *testing.T) {
		// Mock GetUser with record found
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{}, nil
		})
		defer patches.Reset()

		// should return an err
		_, err := svc.Register(context.Background(), request)
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
	request := &proto.LoginRequest{
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
		patches.ApplyFunc(hashing.CheckPasswordHashBcrypt, func(_, _ string) bool {
			return true
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface{}, _ time.Duration) error {
			return nil
		})

		// Test and compare
		resp, err := svc.Login(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_Login got error %v", err)
		}

		if resp.GetSessionId() == "" {
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
		_, err := svc.Login(context.Background(), request)
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
		patches.ApplyFunc(hashing.CheckPasswordHashBcrypt, func(_, _ string) bool {
			return false
		})
		// Test and compare
		_, err := svc.Login(context.Background(), request)
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
		patches.ApplyFunc(hashing.CheckPasswordHashBcrypt, func(_, _ string) bool {
			return true
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface{}, _ time.Duration) error {
			return errors.New("error")
		})

		// Test and compare
		_, err := svc.Login(context.Background(), request)
		if err == nil {
			t.Errorf("TestUserService_Login should return err but didn't")
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	nickname := "test_nickname"
	profilePic := "test_profile_url"
	sessionId := "test_session_id"

	// Input
	request := &proto.GetUserRequest{
		SessionId: sessionId,
	}

	t.Run("normal getUser from cache", func(t *testing.T) {
		want := &proto.GetUserReply{
			Username:   username,
			Nickname:   nickname,
			ProfilePic: profilePic,
		}
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "UserAuth", func(_ UserService, _ string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(_ *dao.RedisCache, _ context.Context, _ string) (string, error) {
			v, _ := json.Marshal(want)
			return string(v), nil
		})

		// Test and compare
		resp, err := svc.GetUser(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}

		if want.Nickname != resp.GetNickname() || want.Username != resp.GetUsername() || want.ProfilePic != resp.GetProfilePic() {
			t.Errorf("TestUserService_GetUser want %v got %v", want, resp)
		}
	})

	t.Run("normal getUser from db", func(t *testing.T) {
		want := &proto.GetUserReply{
			Username:   username,
			Nickname:   nickname,
			ProfilePic: profilePic,
		}
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "UserAuth", func(_ UserService, _ string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(_ *dao.RedisCache, _ context.Context, _ string) (string, error) {
			return "", errors.New("error")
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				Name:       username,
				Nickname:   nickname,
				ProfilePic: profilePic,
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface{}, _ time.Duration) error {
			return nil
		})

		// Test and compare
		resp, err := svc.GetUser(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_GetUser got error %v", err)
		}
		if want.Nickname != resp.GetNickname() || want.Username != resp.GetUsername() || want.ProfilePic != resp.GetProfilePic() {
			t.Errorf("TestUserService_GetUser want %v got %v", want, resp)
		}
	})
}

func TestUserService_EditUser(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	var userId uint32 = 0
	username := "test_username"
	nickname := "test_nickname"
	profilePic := "test_profile_url"
	sessionId := "test_session_id"

	// Input
	request := &proto.EditUserRequest{
		SessionId:  sessionId,
		Nickname:   nickname,
		ProfilePic: profilePic,
	}

	t.Run("normal edit user", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "UserAuth", func(_ UserService, _ string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				CommonModel: &models.CommonModel{
					ID: userId,
				},
				Name:       username,
				Nickname:   nickname,
				ProfilePic: profilePic,
				Status:     uint8(constant.EnabledStatus),
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "UpdateUser", func(_ *dao.Dao, _ uint32, _, _ string) error {
			return nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc), "UpdateUserCache", func(_ UserService, _ string) error {
			return nil
		})

		// Test and compare
		_, err := svc.EditUser(context.Background(), request)
		if err != nil {
			t.Errorf("TestUserService_EditUser got error %v", err)
		}
	})
	t.Run("update failed", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc), "UserAuth", func(_ UserService, _ string) (string, error) {
			return username, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				CommonModel: &models.CommonModel{
					ID: userId,
				},
				Name:       username,
				Nickname:   nickname,
				ProfilePic: profilePic,
				Status:     uint8(constant.EnabledStatus),
			}, nil
		})
		patches.ApplyMethod(reflect.TypeOf(svc.dao), "UpdateUser", func(_ *dao.Dao, _ uint32, _, _ string) error {
			return errors.New("error")
		})

		// Test and compare
		_, err := svc.EditUser(context.Background(), request)
		if err == nil {
			t.Error("TestUserService_EditUser should return error but didn't")
		}
	})
}

func TestUserService_UserAuth(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	sessionId := "test_session_id"

	t.Run("normal user auth", func(t *testing.T) {
		want := username
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(_ *dao.RedisCache, _ context.Context, _ string) (string, error) {
			return username, nil
		})
		defer patches.Reset()

		// Test and compare
		resp, err := svc.UserAuth(sessionId)
		if err != nil {
			t.Errorf("TestUserService_UserAuth got error %v", err)
		}
		if want != resp {
			t.Errorf("TestUserService_UserAuth want %v got %v", want, resp)
		}
	})
	t.Run("user auth failed", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.cache), "Get", func(_ *dao.RedisCache, _ context.Context, _ string) (string, error) {
			return "", errors.New("error")
		})
		defer patches.Reset()

		// Test and compare
		_, err := svc.UserAuth(sessionId)
		if err == nil {
			t.Errorf("TestUserService_EditUser should return error but didn't")
		}
	})
}

func TestUserService_UpdateUserCache(t *testing.T) {
	svc := NewUserService(context.Background())

	// Mock stuffs
	username := "test_username"
	password := "test_password"
	nickname := "test_nickname"

	t.Run("normal update user cache", func(t *testing.T) {
		// Mock DAO call
		patches := gomonkey.ApplyMethod(reflect.TypeOf(svc.dao), "GetUserByName", func(_ *dao.Dao, _ string) (models.UserTab, error) {
			return models.UserTab{
				Name:     username,
				Nickname: nickname,
				Password: password, // It's hashed actually.
			}, nil
		})
		defer patches.Reset()
		patches.ApplyMethod(reflect.TypeOf(svc.cache), "Set", func(_ *dao.RedisCache, _ context.Context, _ string, _ interface{}, _ time.Duration) error {
			return nil
		})

		// Test and compare
		err := svc.UpdateUserCache(username)
		if err != nil {
			t.Errorf("TestUserService_UpdateUserCache got error %v", err)
		}
	})
}
