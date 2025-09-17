package domain

import "gochat/internal/utils"

type Response struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data,omitempty"`
}

func NewResponse(code ResCode, msg string, data ...any) *Response {
	res := &Response{
		Code: code,
		Msg:  msg,
		Data: utils.NormalizeVarArgs(data...),
	}

	return res
}

func NewDefaultResponse(code ResCode, data ...any) *Response {
	return NewResponse(code, code.Msg(), data...)
}

func NewSuccessResponse(data ...any) *Response {
	return NewDefaultResponse(CodeSuccess, data...)
}

type SignupResponse struct {
	UserNumber UserNumber `json:"user_number"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateRoomResponse struct {
	RoomNumber  RoomNumber `json:"room_number,string"`
	RoomName    string     `json:"room_name"`
	Owner       UserNumber `json:"owner,string"`
	Description string     `json:"description"`
	MaxUsers    int        `json:"max_users"`
}

var DefaultResponse = NewSuccessResponse()

// middleware || handler response
var (
	TimeoutResponse      = NewDefaultResponse(CodeTimeout)
	UnauthorizedResponse = NewDefaultResponse(CodeUnauthorized)
	InvalidTokenResponse = NewDefaultResponse(CodeInvalidToken)
)

// usecase logic response
var (
	UserExistResponse     = NewDefaultResponse(CodeUserExist)
	UserNotExistResponse  = NewDefaultResponse(CodeUserNotExist)
	WrongPasswordResponse = NewDefaultResponse(CodeWrongPassword)
	RoomExistResponse     = NewDefaultResponse(CodeRoomExist)
	RoomNotExistResponse  = NewDefaultResponse(CodeRoomNotExist)
	RoomIsFullResponse    = NewDefaultResponse(CodeRoomIsFull)
	WrongSecretResponse   = NewDefaultResponse(CodeWrongSecret)
	HasJoinedResponse     = NewDefaultResponse(CodeHasJoined)
	NotJoinedResponse     = NewDefaultResponse(CodeNotJoined)
)
