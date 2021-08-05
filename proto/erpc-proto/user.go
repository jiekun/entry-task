// @Author: 2014BDuck
// @Date: 2021/8/5

package erpc_proto

type UserService interface {
	Login(*LoginRequest) (*LoginReply, error)
	Register(*RegisterRequest) (*RegisterReply, error)
	GetUser(*GetUserRequest) (*GetUserReply, error)
	EditUser(*EditUserRequest) (*EditUserReply, error)
}

type LoginRequest struct {
	Username string
	Password string
}

type RegisterRequest struct {
	Username  string
	Password  string
	Nickname  string
	ProfilePic string
}
type EditUserRequest struct {
	SessionId  string
	Nickname   string
	ProfilePic string
}
type GetUserRequest struct {
	SessionId string
}
type LoginReply struct {
	SessionId string
}

type RegisterReply struct{}

type EditUserReply struct{}

type GetUserReply struct {
	Username   string
	Nickname   string
	ProfilePic string
}
