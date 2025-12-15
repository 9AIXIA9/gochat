package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/shared/api"

	"go.uber.org/zap"
)

type Request struct {
	Topic   Topic           `json:"topic" validate:"required"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Response struct {
	Topic   Topic `json:"topic"`
	Payload any   `json:"payload,omitempty"`
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

func (r *Router) Route(ctx context.Context, message []byte) *Response {
	request := new(Request)
	if err := json.Unmarshal(message, request); err != nil {
		return &Response{
			Payload: api.NewResponse(api.CodeInvalidParam),
		}
	}

	msg, err := r.validator.Validate(ctx, request)
	if err != nil {
		return &Response{
			Topic:   request.Topic,
			Payload: api.NewResponse(api.CodeServerError),
		}
	}

	if len(msg) != 0 {
		return &Response{
			Topic:   request.Topic,
			Payload: api.NewResponseWithMessage(api.CodeInvalidParam, msg),
		}
	}

	h, ok := r.handlers[request.Topic]
	if !ok {
		if r.notFound != nil {
			h = r.notFound
		} else {
			zap.L().Warn("websocket: topic is not found", zap.String("topic", request.Topic.String()))
			return &Response{
				Topic:   request.Topic,
				Payload: api.NewResponse(api.CodeNotFound),
			}
		}
	}

	return &Response{
		Topic:   request.Topic,
		Payload: h.Handle(SetTopic(ctx, request.Topic), request.Payload),
	}
}
