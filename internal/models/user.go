// @Author: 2014BDuck
// @Date: 2021/5/16

package models

import (
	"github.com/2014bduck/entry-task/internal/constant"
	"gorm.io/gorm"
)

type UserTab struct {
	*CommonModel
	Name      string `json:"name"`       // 用户名
	Nickname  string `json:"nickname"`   // 用户昵称
	AvatarUrl string `json:"avatar_url"` // 用户头像路径
	Password  string `json:"password"`   // 用户登陆密码（加盐）
	Status    uint8  `json:"status"`     // 用户状态 0-enabled 1-disabled
}

func (u UserTab) TableName() string {
	return "user_tab"
}

type UserSessionTab struct {
	*CommonModel
	SessionID  string `json:"session_id"`  // 用户SessionID
	UserID     uint32 `json:"user_id"`     // 用户ID
	ExpireTime uint32 `json:"expire_time"` // 用户Session失效时间
}

func (us UserSessionTab) TableName() string {
	return "user_session_tab"
}

func (u UserTab) GetValidUser(db *gorm.DB) (UserTab, error) {
	var userTab UserTab
	db = db.Where("name = ? AND password = ? AND status = ?", u.Name, u.Password, constant.EnabledStatus)
	err := db.First(&userTab).Error
	if err != nil {
		return userTab, err
	}
	return userTab, nil
}
