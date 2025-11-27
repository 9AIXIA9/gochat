package http

import "net/http"

// BusinessCode 为HTTP码进行补充说明
type BusinessCode int

const CodeSuccess BusinessCode = 200

const CodeServerError BusinessCode = 500

// 业务错误
const (
	CodeTimeout BusinessCode = 400 + iota
	CodeInvalidParam
	CodeInvalidToken
	CodeNotBelongTo
	CodeMaxReached
	CodeInvalidRoom
	CodeOwnerCantLeave
	CodeNotFound
)

func (c BusinessCode) ToHTTPCode() int {
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
	case CodeNotBelongTo:
		return http.StatusConflict
	case CodeInvalidRoom:
		return http.StatusNotAcceptable
	case CodeMaxReached:
		return http.StatusBadRequest
	case CodeOwnerCantLeave:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

func (c BusinessCode) String() string {
	switch c {
	case CodeSuccess:
		return "success"
	case CodeInvalidParam:
		return "invalid param"
	case CodeInvalidToken:
		return "invalid token"
	case CodeTimeout:
		return "request timeout"
	case CodeNotBelongTo:
		return "not belong to"
	case CodeMaxReached:
		return "reach max"
	case CodeInvalidRoom:
		return "invalid room"
	case CodeOwnerCantLeave:
		return "owner cannot leave the room"
	case CodeNotFound:
		return "not found"
	default:
		return "unknown error"
	}
}
