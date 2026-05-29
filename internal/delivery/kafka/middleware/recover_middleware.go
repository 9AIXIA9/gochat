package middleware

import (
	"context"
	"fmt"
	"gochat/internal/infrastructure/kafka"
	"gochat/pkg/ctxutil"
	"runtime/debug"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func NewRecoverMiddleware() kafka.Middleware {
	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) (err error) {
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					zap.L().Error("panic recovered",
						zap.Any("panic", r),
						zap.String("stack", stack),
						zap.String("user_id", ctxutil.UserIDFrom(ctx).String()),
						zap.String("topic", *message.TopicPartition.Topic),
					)
					// Return a safe, generic error to the client
					err = fmt.Errorf("panic recovered: %v", r)
				}
			}()

			return next.Handle(ctx, message)
		})
	}
}
