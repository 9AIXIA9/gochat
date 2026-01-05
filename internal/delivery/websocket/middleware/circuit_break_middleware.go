package middleware

import (
	"context"
	"errors"
	"fmt"
	"gochat/pkg/ctxutil"

	"gochat/internal/infrastructure/breaker"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"

	"github.com/sony/gobreaker/v2"
)

func NewCircuitBreakMiddleware(conf *breaker.Config) websocket.Middleware {
	b := breaker.NewBreaker[any](conf)
	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
			var resp *api.Response

			_, err := b.Execute(func() (any, error) {
				r := next.Handle(ctx, data)
				resp = r

				ctxErr := ctxutil.ErrorFrom(ctx)
				if ctxErr != nil {
					return nil, fmt.Errorf("ctx error: %w", ctxErr)
				}
				if r == nil {
					return nil, fmt.Errorf("nil response from handler")
				}
				if r.Code == api.CodeServerError {
					return nil, fmt.Errorf("server error code: %d", r.Code)
				}
				return nil, nil
			})

			if err != nil && (errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests)) {
				return api.NewResponse(api.CodeServiceUnavailable)
			}
			return resp
		})
	}
}
