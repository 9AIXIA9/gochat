package domain

import "net/http"

type ResCode int

// 定义错误码常量

// CodeSuccess 成功响应
const CodeSuccess ResCode = 200

// CodeServerBusy 服务端错误
const CodeServerBusy ResCode = 500

// 客户端错误 4xx
const (
	CodeInvalidParam ResCode = 400 + iota
	CodeUnauthorized         //refresh token
	CodeInvalidToken         //auth token
	CodeUserExist
	CodeUserNotExist
	CodeWrongPassword
	CodeRoomNotExist
	CodeRoomExist
	CodeRoomIsFull
	CodeWrongSecret
	CodeHasJoined
	CodeNotJoined
	CodeTimeout
)

func (c ResCode) ToHTTP() int {
	// 优先对已定义的业务码做精确映射
	switch c {
	case CodeSuccess:
		return http.StatusOK
	case CodeTimeout:
		return http.StatusRequestTimeout
	case CodeInvalidParam:
		return http.StatusBadRequest
	case CodeUnauthorized, CodeInvalidToken, CodeWrongPassword:
		return http.StatusUnauthorized
	case CodeUserExist, CodeRoomExist, CodeHasJoined:
		return http.StatusConflict
	case CodeUserNotExist, CodeRoomNotExist:
		return http.StatusNotFound
	case CodeRoomIsFull, CodeNotJoined, CodeWrongSecret:
		return http.StatusForbidden
	case CodeServerBusy:
		// 服务繁忙适合返回 503
		return http.StatusServiceUnavailable
	}

	// 未显式列出的码按范围回退
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
	case CodeInvalidToken:
		return "无效的token"
	case CodeUserExist:
		return "用户已存在"
	case CodeUserNotExist:
		return "用户不存在"
	case CodeWrongPassword:
		return "用户名或密码错误"
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
