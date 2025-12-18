package websocket

import (
	"context"
	"encoding/json"
	"gochat/internal/shared/api"

	"go.uber.org/zap"
)

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

func (r *Router) Route(ctx context.Context, payload []byte) *Message {
	request := new(Request)
	if err := json.Unmarshal(payload, request); err != nil {
		return &Message{
			Body: api.NewResponse(api.CodeInvalidParam),
		}
	}

	msg, err := r.validator.Validate(ctx, request)
	if err != nil {
		return &Message{
			Topic: request.Topic,
			Body:  api.NewResponse(api.CodeServerError),
		}
	}

	if len(msg) != 0 {
		return &Message{
			Topic: request.Topic,
			Body:  api.NewResponseWithMessage(api.CodeInvalidParam, msg),
		}
	}

	h, ok := r.handlers[request.Topic]
	if !ok {
		if r.notFound != nil {
			h = r.notFound
		} else {
			zap.L().Warn("websocket: topic is not found", zap.String("topic", request.Topic.String()))
			return &Message{
				Topic: request.Topic,
				Body:  api.NewResponse(api.CodeNotFound),
			}
		}
	}

	return &Message{
		Topic: request.Topic,
		Body:  h.Handle(SetTopic(ctx, request.Topic), request.Payload),
	}
}
