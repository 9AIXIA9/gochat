package websocket

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

//TODO为处理器添加信息 链路和指标

type Router struct {
	handlers    map[RequestTopic]Handler
	middlewares []Middleware
	notFound    Handler
}

func NewRouter() *Router {
	return &Router{
		handlers:    make(map[RequestTopic]Handler),
		middlewares: make([]Middleware, 0),
		notFound: HandlerFunc(func(ctx context.Context, req *Request) (*Response, error) {
			return nil, fmt.Errorf("no handler for topic %s", req.RequestTopic)
		}),
	}
}

func (r *Router) Use(mw ...Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

func (r *Router) Handle(topic RequestTopic, h Handler) {
	if _, ok := r.handlers[topic]; ok {
		zap.L().Warn("handler already registered for topic", zap.String("topic", string(topic)))
		r.handlers[topic] = h
		return
	}
	if len(r.middlewares) > 0 {
		h = Chain(h, r.middlewares...)
	}
	r.handlers[topic] = h
}

func (r *Router) Route(ctx context.Context, req *Request) (*Response, error) {
	h, ok := r.handlers[req.RequestTopic]
	if !ok {
		h = r.notFound
	}
	resp, err := h.Handle(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
