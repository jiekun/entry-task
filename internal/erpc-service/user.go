// @Author: 2014BDuck
// @Date: 2021/8/3

package erpc_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/2014bduck/entry-task/global"
	"github.com/2014bduck/entry-task/internal/constant"
	"github.com/2014bduck/entry-task/internal/dao"
	"github.com/2014bduck/entry-task/pkg/hashing"
	"github.com/2014bduck/entry-task/pkg/rpc/erpc"
	proto "github.com/2014bduck/entry-task/proto/erpc-proto"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	ctx   context.Context
	dao   *dao.Dao
	cache *dao.RedisCache
}

func NewUserService(ctx context.Context) UserService {
	svc := UserService{ctx: ctx}
	svc.dao = dao.New(global.DBEngine)
	svc.cache = dao.NewCache(global.CacheClient)

	return svc
}

func (svc UserService) RegisterUserService(s *erpc.Server) {
	s.Register("Login", svc.Login, proto.LoginRequest{}, proto.LoginReply{})
	s.Register("Register", svc.Register, proto.RegisterRequest{}, proto.RegisterReply{})
	s.Register("EditUser", svc.EditUser, proto.EditUserRequest{}, proto.EditUserReply{})
	s.Register("GetUser", svc.GetUser, proto.GetUserRequest{}, proto.GetUserReply{})
}

func (svc UserService) Login(r proto.LoginRequest) (*proto.LoginReply, error) {
	// Implement distributed lock with Redis if necessary
	// Key, Value := xxx, xxx
	// if Redis.SetNX(Key, Value, ttl){ Do Business Logic }

	// Find user
	user, err := svc.dao.GetUserByName(r.Username)
	if err != nil {
		return nil, err
	}

	// Invalid cases
	hashedPass := hashing.HashPassword(r.Password)
	if user.Password != hashedPass {
		return nil, errors.New("svc.UserLogin: pwd incorrect")
	} else if user.Status != uint8(constant.EnabledStatus) {
		return nil, errors.New("svc.UserLogin: status disabled")
	}

	// Validation success
	// Setting session cache
	sessionID := uuid.NewV4()
	err = svc.cache.Cache.Set(svc.ctx, constant.SessionIDCachePrefix+sessionID.String(), []byte(r.Username), 0).Err()

	if err != nil {
		return nil, err
	}
	return &proto.LoginReply{SessionId: sessionID.String()}, nil
}

func (svc UserService) Register(r *proto.RegisterRequest) (*proto.RegisterReply, error) {
	// Validate username if existed
	_, err := svc.dao.GetUserByName(r.Username)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("svc.UserRegister: username existed")
	}

	// Add Salt to pass
	hashedPass := hashing.HashPassword(r.Password)

	// Create User to DB
	_, err = svc.dao.CreateUser(r.Username, hashedPass, r.Nickname, r.ProfilePic, uint8(constant.EnabledStatus))
	if err != nil {
		return nil, fmt.Errorf("svc.UserRegister: CreateUser error: %v", err)
	}

	return &proto.RegisterReply{}, nil
}

func (svc UserService) EditUser(r *proto.EditUserRequest) (*proto.EditUserReply, error) {
	// Get Username
	username, err := svc.UserAuth(r.SessionId)
	if err != nil {
		return nil, err
	}

	// Query current user
	user, err := svc.dao.GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("svc.UserEdit: %v", err)
	}

	// Validate user status
	if constant.Status(user.Status) != constant.EnabledStatus {
		return nil, errors.New("svc.UserEdit: Invalid user status")
	}

	// Update user data
	err = svc.dao.UpdateUser(user.ID, r.Nickname, r.ProfilePic)
	if err != nil {
		return nil, fmt.Errorf("svc.UserEdit: %v", err)
	}

	// Update Cache
	_ = svc.UpdateUserCache(username)

	return &proto.EditUserReply{}, nil
}

func (svc UserService) GetUser(r *proto.GetUserRequest) (*proto.GetUserReply, error) {
	// Get Username
	username, err := svc.UserAuth(r.SessionId)
	if err != nil {
		return nil, err
	}

	cacheKey := constant.UserProfileCachePrefix + username

	// Try loading user info from cache
	userProfCache, err := svc.cache.Cache.Get(svc.ctx, cacheKey).Result()
	if err == nil {
		userGetCacheResp := proto.GetUserReply{}
		err = json.Unmarshal([]byte(userProfCache), &userGetCacheResp)
		if err != nil {
			global.Logger.Errorf("svc.UserGet: Unmarshal cache failed: %v", err)
		} else {
			return &userGetCacheResp, nil
		}
	}

	// Query user from DB
	user, err := svc.dao.GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("svc.UserGet: %v", err)
	}
	userGetResp := &proto.GetUserReply{
		Username:   user.Name,
		Nickname:   user.Nickname,
		ProfilePic: user.ProfilePic,
	}

	// Set user to cache
	cacheUser, _ := json.Marshal(userGetResp)
	err = svc.cache.Cache.Set(svc.ctx, cacheKey, cacheUser, 3600*24*time.Second).Err() // Omit error
	if err != nil {
		global.Logger.Errorf("svc.UserGet: set cache failed: %v", err)
	}

	return userGetResp, nil
}

func (svc UserService) UpdateUserCache(username string) error {
	cacheKey := constant.UserProfileCachePrefix + username

	// Query user from DB
	user, err := svc.dao.GetUserByName(username)
	if err != nil {
		return fmt.Errorf("svc.UserGet: %v", err)
	}
	userGetResp := &proto.GetUserReply{
		Username:   user.Name,
		Nickname:   user.Nickname,
		ProfilePic: user.ProfilePic,
	}

	// Set user to cache
	cacheUser, _ := json.Marshal(userGetResp)
	err = svc.cache.Cache.Set(svc.ctx, cacheKey, cacheUser, 3600*24*time.Second).Err() // Omit error
	if err != nil {
		global.Logger.Errorf("svc.UserGet: set cache failed: %v", err)
	}
	return nil
}

func (svc UploadService) UserAuth(sessionID string) (string, error) {
	username, err := svc.cache.Cache.Get(svc.ctx, constant.SessionIDCachePrefix+sessionID).Result()

	if err != nil {
		return "", errors.New("svc.UserAuth failed")
	}
	return string(username), nil
}
