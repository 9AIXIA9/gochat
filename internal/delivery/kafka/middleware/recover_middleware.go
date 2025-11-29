package middleware

import (
	"context"
	"gochat/internal/infrastructure/kafka"
	myErrors "gochat/internal/shared/errors"
	"gochat/pkg/utils"
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
						zap.String("user_id", utils.GetUserID(ctx).String()),
						zap.String("topic", utils.GetWebsocketTopic(ctx)),
					)
					// Return a safe, generic error to the client
					err = myErrors.ErrInternalError
				}
			}()

			return next.Handle(ctx, message)
		})
	}
}
