package api

import "net/http"

type Code int

// 业务错误
const (
	CodeSuccess       Code = 200
	CodeBusinessError Code = 300
	CodeServerError   Code = 500
	CodeTimeout       Code = 400 + iota
	CodeInvalidParam
	CodeInvalidToken
	CodeNotFound
	CodeServiceUnavailable
)

func (c Code) ToHTTPCode() int {
	switch c {
	case CodeSuccess:
		return http.StatusOK
	case CodeBusinessError:
		return http.StatusBadRequest
	case CodeServerError:
		return http.StatusInternalServerError
	case CodeTimeout:
		return http.StatusRequestTimeout
	case CodeInvalidParam:
		return http.StatusBadRequest
	case CodeInvalidToken:
		return http.StatusUnauthorized
	case CodeNotFound:
		return http.StatusNotFound
	case CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func (c Code) DefaultMessage() string {
	switch c {
	case CodeSuccess:
		return "success"
	case CodeBusinessError:
		return "business error"
	case CodeServerError:
		return "server error"
	case CodeTimeout:
		return "request timeout"
	case CodeInvalidParam:
		return "invalid parameter"
	case CodeInvalidToken:
		return "invalid token"
	case CodeNotFound:
		return "not found"
	case CodeServiceUnavailable:
		return "service unavailable"
	default:
		return "unknown error"
	}
}
