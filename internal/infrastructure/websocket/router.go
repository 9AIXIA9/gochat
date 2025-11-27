package websocket

import (
	"context"
	"encoding/json"
	"gochat/pkg/utils"

	"go.uber.org/zap"
)

type Topic string

func (t Topic) String() string {
	return string(t)
}

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
	middlewares []Middleware
	notFound    Handler
}

func NewRouter() *Router {
	return &Router{
		handlers:    make(map[Topic]Handler),
		middlewares: make([]Middleware, 0),
		notFound:    nil,
	}
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) NoRoute(h Handler, middlewares ...Middleware) {
	wrapped := chainHandlers(h, middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.notFound = wrapped
}

func (r *Router) Handle(topic Topic, h Handler, middlewares ...Middleware) {
	// Apply route-level middlewares first, then global middlewares
	wrapped := chainHandlers(h, middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.handlers[topic] = wrapped
}

func (r *Router) Route(ctx context.Context, request *Request) *Response {
	h, ok := r.handlers[request.Topic]
	if !ok {
		if r.notFound != nil {
			h = r.notFound
		} else {
			zap.L().Debug("websocket: topic is not found", zap.String("topic", request.Topic.String()))
			data, mErr := json.Marshal(&ErrorData{Message: "topic is not found"})
			if mErr != nil {
				return nil
			}
			return &Response{
				Topic: request.Topic,
				Data:  data,
			}
		}
	}

	ctx = utils.SetWebsocketTopic(ctx, request.Topic.String())

	resp, err := h.Handle(ctx, request.Data)
	if err != nil {
		data, mErr := json.Marshal(&ErrorData{Message: err.Error()})
		if mErr != nil {
			return nil
		}
		resp = data
	}
	return &Response{
		Topic: request.Topic,
		Data:  resp,
	}
}
