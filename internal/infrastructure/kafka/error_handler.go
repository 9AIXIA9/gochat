package kafka

import (
	"context"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error, message *ckafka.Message)
}

var _ ErrorHandler = ErrorHandlerFunc(nil)

type ErrorMiddleware func(next ErrorHandler) ErrorHandler

type ErrorHandlerFunc func(ctx context.Context, err error, message *ckafka.Message)

func (f ErrorHandlerFunc) Handle(ctx context.Context, err error, message *ckafka.Message) {
	f(ctx, err, message)
}

func chainErrorHandlers(mainErrorHandler ErrorHandler, middlewares []ErrorMiddleware) ErrorHandler {
	if len(middlewares) == 0 {
		return mainErrorHandler
	}

	// Apply in reverse so the first middleware becomes the outermost wrapper.
	wrapped := mainErrorHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
