package domain

import "net/http"

type ResCode int

// 定义错误码常量
const (
	CodeSuccess ResCode = 200

	// 客户端错误 4xx
	CodeInvalidParam ResCode = 400 + iota
	CodeUnauthorized
	CodeUserExist
	CodeUserNotExist
	CodeWrongPassword
	CodeInvalidToken
	CodeRoomNotExist
	CodeRoomExist
	CodeRoomIsFull
	CodeWrongSecret
	CodeHasJoined
	CodeNotJoined
	CodeTimeout

	// 服务端错误 5xx
	CodeServerBusy ResCode = 500 + iota
)

func (c ResCode) ToHTTP() int {
	switch {
	case c >= 200 && c < 300:
		return http.StatusOK
	case c == CodeTimeout:
		return http.StatusRequestTimeout
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
	case CodeRoomExist:
		return "房间已存在"
	case CodeRoomIsFull:
		return "房间已满"
	case CodeServerBusy:
		return "服务繁忙"
	case CodeWrongSecret:
		return "房间号或密钥错误"
	case CodeHasJoined:
		return "已加入房间"
	case CodeNotJoined:
		return "未加入房间"
	case CodeTimeout:
		return "请求超时"
	default:
		return "未知错误"
	}
}
