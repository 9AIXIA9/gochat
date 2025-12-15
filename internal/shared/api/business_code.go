package api

import "net/http"

type Code int

const CodeSuccess Code = 200

const CodeServerError Code = 500

// 业务错误
const (
	CodeTimeout Code = 400 + iota
	CodeInvalidParam
	CodeInvalidToken
	CodeNotFound
)

func (c Code) ToHTTPCode() int {
	switch c {
	case CodeSuccess:
		return http.StatusOK
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

	default:
		return http.StatusInternalServerError
	}
}

func (c Code) String() string {
	switch c {
	case CodeSuccess:
		return "success"
	case CodeInvalidParam:
		return "invalid param"
	case CodeInvalidToken:
		return "invalid token"
	case CodeTimeout:
		return "request timeout"
	case CodeNotFound:
		return "not found"
	default:
		return "unknown error"
	}
}
