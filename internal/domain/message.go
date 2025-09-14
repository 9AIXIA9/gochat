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
