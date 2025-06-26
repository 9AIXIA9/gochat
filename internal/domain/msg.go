package domain

import "net/http"

type ResCode int

// 定义错误码常量
const (
	CodeSuccess ResCode = 200

	// 客户端错误 4xx
	CodeInvalidParam  ResCode = 400
	CodeUnauthorized  ResCode = 401
	CodeForbidden     ResCode = 403
	CodeNotFound      ResCode = 404
	CodeUserExist     ResCode = 410
	CodeUserNotExist  ResCode = 411
	CodeWrongPassword ResCode = 412
	CodeInvalidToken  ResCode = 413
	CodeRoomNotExist  ResCode = 414
	CodeRoomIsFull    ResCode = 415

	// 服务端错误 5xx
	CodeServerBusy     ResCode = 500
	CodeInvalidRequest ResCode = 501
)

type Message struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data,omitempty"`
}

func (c ResCode) ToHTTP() int {
	switch {
	case c >= 200 && c < 300:
		return http.StatusOK
	case c >= 400 && c < 500:
		return http.StatusBadRequest
	case c >= 500:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

func (c ResCode) Msg() string {
	switch c {
	case CodeSuccess:
		return "成功"
	case CodeInvalidParam:
		return "请求参数错误"
	case CodeUnauthorized:
		return "未认证"
	case CodeForbidden:
		return "没有权限"
	case CodeNotFound:
		return "请求资源不存在"
	case CodeUserExist:
		return "用户已存在"
	case CodeUserNotExist:
		return "用户不存在"
	case CodeWrongPassword:
		return "用户名或密码错误"
	case CodeInvalidToken:
		return "无效的token"
	case CodeRoomNotExist:
		return "房间不存在"
	case CodeRoomIsFull:
		return "房间已满"
	case CodeServerBusy:
		return "服务繁忙"
	case CodeInvalidRequest:
		return "请求错误"
	default:
		return "未知错误"
	}
}
