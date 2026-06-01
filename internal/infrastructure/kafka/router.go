package kafka

import (
	"context"
	"gochat/internal/shared/command"
	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Router struct {
	handlers        map[string]Handler
	middlewares     []Middleware
	notFoundHandler Handler
}

func NewRouter() *Router {
	return &Router{
		handlers:        make(map[string]Handler),
		middlewares:     make([]Middleware, 0),
		notFoundHandler: nil,
	}
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) NoRoute(h Handler, middlewares ...Middleware) {
	wrapped := chainHandlers(h, middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.notFoundHandler = wrapped
}

func (r *Router) EventHandle(topic event.Topic, h event.Handler, middlewares ...Middleware) {
	// Apply route-level middlewares first, then global middlewares
	wrapped := chainHandlers(WrapEventHandler(h), middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.handlers[topic.String()] = wrapped
}

func (r *Router) CommandHandle(action command.Action, h command.Handler, middlewares ...Middleware) {
	// Apply route-level middlewares first, then global middlewares
	wrapped := chainHandlers(WrapCommandHandler(h), middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.handlers[action.String()] = wrapped
}

func (r *Router) ReceiptHandle(topic string, h command.ReceiptHandler, middlewares ...Middleware) {
	// Apply route-level middlewares first, then global middlewares
	wrapped := chainHandlers(WrapReceiptHandler(h), middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.handlers[topic] = wrapped
}

func (r *Router) Route(ctx context.Context, message *ckafka.Message) error {
	topic := *message.TopicPartition.Topic
	h, ok := r.handlers[topic]
	if !ok {
		if r.notFoundHandler != nil {
			h = r.notFoundHandler
		} else {
			zap.L().Debug("event: topic is not found", zap.String("topic", topic))
		}
	}

	return h.Handle(ctx, message)
}

func (r *Router) Topics() []string {
	topics := make([]string, 0, len(r.handlers))
	for topic := range r.handlers {
		topics = append(topics, topic)
	}
	return topics
}
