package websocket

import (
	"context"
	"encoding/json"
)

type RequestTopic string

type Request struct {
	RequestTopic RequestTopic    `json:"topic"`
	Data         json.RawMessage `json:"data,omitempty"`
}
type ResponseTopic string

type Response struct {
	ResponseTopic ResponseTopic   `json:"topic"`
	Data          json.RawMessage `json:"data,omitempty"`
}

type Handler interface {
	Handle(ctx context.Context, request *Request) (*Response, error)
}

type HandlerFunc func(ctx context.Context, request *Request) (*Response, error)

func (f HandlerFunc) Handle(ctx context.Context, request *Request) (*Response, error) {
	return f(ctx, request)
}

type Middleware func(Handler) Handler

func Chain(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

//func MiddlewareFunc(f func(ctx context.Context, req *Request, next Handler) (*Response, error)) Middleware {
//	return func(next Handler) Handler {
//		return HandlerFunc(func(ctx context.Context, req *Request) (*Response, error) {
//			return f(ctx, req, next)
//		})
//	}
//}
