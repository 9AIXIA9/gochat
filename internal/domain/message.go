package domain

import "gochat/internal/utils"

type Message struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data any     `json:"data,omitempty"`
}

func NewMessage(code ResCode, msg string, data ...any) *Message {
	res := &Message{
		Code: code,
		Msg:  msg,
		Data: utils.NormalizeVarArgs(data...),
	}

	return res
}

func NewDefaultMessage(code ResCode, data ...any) *Message {
	return NewMessage(code, code.Msg(), data...)
}

func NewSuccessMessage(data ...any) *Message {
	return NewDefaultMessage(CodeSuccess, data...)
}
