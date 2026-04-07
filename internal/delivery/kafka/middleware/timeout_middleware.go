package middleware

import (
	"context"
	"errors"
	"fmt"
	"gochat/internal/infrastructure/kafka"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// NewTimeoutMiddleware enforces a per-message timeout for kafka route handlers.
func NewTimeoutMiddleware(duration time.Duration) kafka.Middleware {
	if duration <= 0 {
		return func(next kafka.Handler) kafka.Handler { return next }
	}

	return func(next kafka.Handler) kafka.Handler {
		return kafka.HandlerFunc(func(ctx context.Context, message *ckafka.Message) error {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, duration)
			defer cancel()

			errCh := make(chan error, 1)
			go func() {
				errCh <- next.Handle(ctxWithTimeout, message)
			}()

			select {
			case err := <-errCh:
				return err
			case <-ctxWithTimeout.Done():
				if errors.Is(ctxWithTimeout.Err(), context.DeadlineExceeded) {
					return fmt.Errorf("kafka handler timeout after %s: %w", duration, ctxWithTimeout.Err())
				}
				return ctxWithTimeout.Err()
			}
		})
	}
}
