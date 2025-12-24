package middleware

import (
	"context"
	"errors"
	"fmt"
	"gochat/internal/infrastructure/breaker"
	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sony/gobreaker/v2"
)

// NewCircuitBreakMiddleware wraps kafka handlers with a circuit breaker to short-circuit on high failure rates.
func NewCircuitBreakMiddleware(conf *breaker.Config) kafka.Middleware {
	if conf == nil {
		conf = &breaker.Config{}
	}
	// Ensure sensible defaults when caller skips validation.
	_ = conf.Validate()

	b := breaker.NewBreaker[any](conf)

	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			_, err := b.Execute(func() (any, error) {
				if err := next.Handle(ctx, message); err != nil {
					return nil, err
				}

				if ctxErr := ctx.Err(); ctxErr != nil {
					return nil, ctxErr
				}

				return nil, nil
			})

			if err == nil {
				return nil
			}

			if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
				return fmt.Errorf("kafka consumer circuit open: %w", err)
			}

			return err
		})
	}
}
