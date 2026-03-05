package api

import (
	"encoding/json"
)

type Response struct {
	Code    Code   `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty" swaggertype:"object"`
}

func (r *Response) String() string {
	data, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(data)
}

func NewResponse(code Code) *Response {
	return &Response{
		Code:    code,
		Message: code.DefaultMessage(),
	}
}

func NewResponseWithMessage(code Code, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
	}
}

func NewResponseWithData(data any) *Response {
	return &Response{
		Code:    CodeSuccess,
		Message: CodeSuccess.DefaultMessage(),
		Data:    data,
	}
}
