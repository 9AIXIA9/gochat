package domain

type Response struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data,omitempty"`
}

func NewResponse(code ResCode, msg string, data ...any) *Response {
	res := &Response{
		Code: code,
		Msg:  msg,
	}

	if len(data) > 0 {
		if l := len(data); l == 1 && data[0] != nil {
			res.Data = data[0]
		} else if l > 1 {
			res.Data = data
		}
	}
	return res
}

func NewResponseWithDefaultMsg(code ResCode, data ...any) *Response {
	return NewResponse(code, code.Msg(), data...)
}

func NewSuccessResponse(data ...any) *Response {
	return NewResponseWithDefaultMsg(CodeSuccess, data...)
}
