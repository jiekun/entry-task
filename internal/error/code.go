// @Author: 2014BDuck
// @Date: 2021/7/11

package errcode

var (
	Success       = NewError(0, "Success")
	ServerError   = NewError(10000000, "Server Error")
	InvalidParams = NewError(10000001, "Invalid Params")
	NotFound      = NewError(10000002, "Not Found")

	ErrorUserLogin    = NewError(20010001, "User Login Failed")
	ErrorUserNotLogin = NewError(20010002, "User Login Required")
	ErrorUserRegister = NewError(20010003, "User Register Failed")
	ErrorUserGet      = NewError(20010004, "User Get Profile Failed")
	ErrorUserEdit     = NewError(20010005, "User Edit Profile Failed")

	ErrorUploadPicFailed = NewError(30010001, "Upload Picture Failed")
)
