package websocket

import (
	"context"
	"encoding/json"
)

//TODO 添加middleware（日志，链路，指标等）

type Topic string

type Request struct {
	Topic Topic           `json:"topic"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type Response struct {
	Topic Topic           `json:"topic"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type Router struct {
	handlers    map[Topic]Handler
	middlewares []Handler
	notFound    Handler
}

func NewRouter(notFound Handler) *Router {
	return &Router{
		handlers:    make(map[Topic]Handler),
		middlewares: make([]Handler, 0),
		notFound:    notFound,
	}
}

func (r *Router) Use(middleware ...Handler) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) Handle(topic Topic, h Handler, middlewares ...Handler) {
	r.handlers[topic] = chainHandlers(h, middlewares)
}

func (r *Router) Route(ctx context.Context, request *Request) *Response {
	h, ok := r.handlers[request.Topic]
	if !ok {
		h = r.notFound
	}

	//TODO 异步处理

	resp, err := h.Handle(ctx, request.Data)
	if err != nil {
		data, err := json.Marshal(&MessageData{Message: err.Error()})
		if err != nil {
			return nil
		}
		resp = data
	}
	return &Response{
		Topic: request.Topic,
		Data:  resp,
	}
}
