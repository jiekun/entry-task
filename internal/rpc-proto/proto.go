// @Author: 2014BDuck
// @Date: 2021/8/3

package rpcproto

type UserLoginRequest struct {
	Username string `form:"username" binding:"required,min=2,max=255"`
	Password string `form:"password" binding:"required,min=2,max=255"`
}

type UserRegisterRequest struct {
	Username   string `form:"username" binding:"required,min=2,max=255"`
	Password   string `form:"password" binding:"required,min=2,max=255"`
	Nickname   string `form:"nickname" binding:"required,min=2,max=255"`
	ProfilePic string `form:"profile_pic" binding:"-"` // Skip validation.
}

type UserEditRequest struct {
	SessionID  string
	Nickname   string `form:"nickname"`
	ProfilePic string `form:"profile_pic"`
}

type UserGetRequest struct {
	SessionID string
}

type UserLoginResponse struct {
	SessionID string `json:"session_id"`
}

type UserRegisterResponse struct{}

type UserEditResponse struct{}

type UserGetResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"`
}
