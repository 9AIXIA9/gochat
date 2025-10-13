package domain

import (
	"encoding/json"
	"gochat/internal/utils"
)

type Response struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data,omitempty"`
}

func (r *Response) MarshalJSON() ([]byte, error) {
	//用匿名结构体来进行默认json序列化 防止无限递归
	type Alias Response
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	// 判断 Data 是否为“空结构体”或 nil
	if utils.IsEmptyData(r.Data) {
		return json.Marshal(&struct {
			Code ResCode `json:"code"`
			Msg  string  `json:"msg"`
		}{
			Code: r.Code,
			Msg:  r.Msg,
		})
	}

	return json.Marshal(aux)
}

func NewResponse(code ResCode, msg string, data ...any) *Response {
	normalizedData := utils.NormalizeVarArgs(data...)

	res := &Response{
		Code: code,
		Msg:  msg,
		Data: normalizedData,
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
	AuthToken    AuthToken    `json:"auth_token"`
	RefreshToken RefreshToken `json:"-"`
}

type RefreshTokenResponse struct {
	AuthToken    AuthToken    `json:"auth_token"`
	RefreshToken RefreshToken `json:"-"`
}

type CreateRoomResponse struct {
	RoomNumber RoomNumber `json:"room_number,string"`
	RoomName   string     `json:"room_name"`
	Owner      UserNumber `json:"owner,string"`
	MaxUsers   int        `json:"max_users"`
}

type SendMessageResponse struct {
	MessageID MessageID `json:"message_id"`
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
	RoomNotExistResponse  = NewDefaultResponse(CodeRoomNotExist)
	RoomIsFullResponse    = NewDefaultResponse(CodeRoomIsFull)
	WrongSecretResponse   = NewDefaultResponse(CodeWrongSecret)
	HasJoinedResponse     = NewDefaultResponse(CodeHasJoined)
	NotJoinedResponse     = NewDefaultResponse(CodeNotJoined)
)
