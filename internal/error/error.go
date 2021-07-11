// @Author: 2014BDuck
// @Date: 2021/7/11

package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"detail"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("Errcode %d existed.", code))
	}
	codes[code] = msg
	return &Error{Code: code, Msg: msg}
}

func (e *Error) GetError() string {
	return fmt.Sprintf("Errcode: %d, errmsg: %s", e.GetCode(), e.GetMsg())
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) GetMsg() string {
	return e.Msg
}

func (e *Error) GetDetails() []string {
	return e.Details
}

func (e *Error) StatusCode() int {
	switch e.GetCode() {
	case Success.GetCode():
		return http.StatusOK
	case ServerError.GetCode():
		return http.StatusInternalServerError
	case InvalidParams.GetCode():
		return http.StatusBadRequest
	case NotFound.GetCode():
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
