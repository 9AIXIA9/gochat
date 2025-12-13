package websocket

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

type Request struct {
	Topic Topic           `json:"topic" validate:"required"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type Response struct {
	Topic Topic           `json:"topic" validate:"required"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type Router struct {
	validator   Validator
	handlers    map[Topic]Handler
	middlewares []Middleware
	notFound    Handler
}

func NewRouter(validator Validator) *Router {
	return &Router{
		validator:   validator,
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
	msg, err := r.validator.Validate(ctx, request)
	if err != nil {
		data, mErr := json.Marshal(&ErrorData{Message: err.Error()})
		if mErr != nil {
			return nil
		}
		return &Response{
			Topic: request.Topic,
			Data:  data,
		}
	}

	if len(msg) != 0 {
		data, mErr := json.Marshal(&ErrorData{
			Message: msg,
		})
		if mErr != nil {
			return nil
		}
		return &Response{
			Topic: request.Topic,
			Data:  data,
		}
	}

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

	ctx = SetTopic(ctx, request.Topic)

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
