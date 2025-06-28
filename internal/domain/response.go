package domain

type Response struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data,omitempty"`
}

func NewResponse(code ResCode, msg string, data ...any) *Response {
	var resData any
	switch len(data) {
	case 0:
		resData = nil
	case 1:
		resData = data[0]
	default:
		resData = data
	}

	return &Response{
		Code: code,
		Msg:  msg,
		Data: resData,
	}
}

func NewResponseWithoutMsg(code ResCode, data ...any) *Response {
	return NewResponse(code, code.Msg(), data)
}

func NewSuccessResponse(data ...any) *Response {
	return NewResponse(CodeSuccess, CodeSuccess.Msg(), data)
}
