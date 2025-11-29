package kafka

import (
	"context"
	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Router struct {
	handlers        map[event.Topic]Handler
	middlewares     []Middleware
	notFoundHandler Handler
}

func NewRouter() *Router {
	return &Router{
		handlers:        make(map[event.Topic]Handler),
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

func (r *Router) Handle(topic event.Topic, h Handler, middlewares ...Middleware) {
	// Apply route-level middlewares first, then global middlewares
	wrapped := chainHandlers(h, middlewares)
	wrapped = chainHandlers(wrapped, r.middlewares)
	r.handlers[topic] = wrapped
}

func (r *Router) Route(ctx context.Context, message *ckafka.Message) error {
	topic := event.Topic(*message.TopicPartition.Topic)
	h, ok := r.handlers[topic]
	if !ok {
		if r.notFoundHandler != nil {
			h = r.notFoundHandler
		} else {
			zap.L().Debug("event: topic is not found", zap.String("topic", topic.String()))
		}
	}

	return h.Handle(ctx, message)
}

func (r *Router) Topics() []string {
	topics := make([]string, 0, len(r.handlers))
	for topic := range r.handlers {
		topics = append(topics, topic.String())
	}
	return topics
}
