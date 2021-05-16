// @Author: 2014BDuck
// @Date: 2021/5/16

package models

type UserTab struct {
	*CommonModel
	Name      string `json:"name"`       // 用户名
	Nickname  string `json:"nickname"`   // 用户昵称
	AvatarUrl string `json:"avatar_url"` // 用户头像路径
	Password  string `json:"password"`   // 用户登陆密码（加盐）
}

func (userTab UserTab) TableName() string {
	return "user_tab"
}

type UserSessionTab struct {
	*CommonModel
	SessionID  string `json:"session_id"`  // 用户SessionID
	UserID     uint32 `json:"user_id"`     // 用户ID
	ExpireTime uint32 `json:"expire_time"` // 用户Session失效时间
}

func (userSessionTab UserSessionTab) TableName() string {
	return "user_session_tab"
}
