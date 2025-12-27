package api

type Response struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty" swaggertype:"object"`
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
