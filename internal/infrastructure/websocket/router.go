package websocket

import (
	"context"
	"fmt"
)

type Router struct {
	handlers    map[RequestTopic]Handler
	middlewares []Middleware
	notFound    Handler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[RequestTopic]Handler),
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
		panic("duplicate handler for topic: " + string(topic))
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
