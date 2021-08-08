// @Author: 2014BDuck
// @Date: 2021/8/5

package erpc_proto

// Function signature for Client. Align with interface.
var (
	Login    func(LoginRequest) (LoginReply, error)
	Register func(RegisterRequest) (RegisterReply, error)
	GetUser  func(GetUserRequest) (GetUserReply, error)
	EditUser func(EditUserRequest) (EditUserReply, error)
)

// UserService interface for Server. Align with signature.
type UserService interface {
	Login(LoginRequest) (*LoginReply, error)
	Register(RegisterRequest) (*RegisterReply, error)
	GetUser(GetUserRequest) (*GetUserReply, error)
	EditUser(EditUserRequest) (*EditUserReply, error)
}

type LoginRequest struct {
	Username string
	Password string
}

type RegisterRequest struct {
	Username   string
	Password   string
	Nickname   string
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
