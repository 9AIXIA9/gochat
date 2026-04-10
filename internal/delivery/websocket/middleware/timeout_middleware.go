package middleware

import (
	"context"
	"errors"
	"gochat/internal/infrastructure/websocket"
	"gochat/internal/shared/api"
	"time"
)

// NewTimeoutMiddleware enforces a per-message timeout for websocket route handlers.
func NewTimeoutMiddleware(duration time.Duration) websocket.Middleware {
	if duration <= 0 {
		return func(next websocket.Handler) websocket.Handler { return next }
	}

	return func(next websocket.Handler) websocket.Handler {
		return websocket.HandlerFunc(func(ctx context.Context, data []byte) *api.Response {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, duration)
			defer cancel()

			respCh := make(chan *api.Response, 1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						respCh <- api.ResponseServerError
					}
				}()
				respCh <- next.Handle(ctxWithTimeout, data)
			}()

			select {
			case resp := <-respCh:
				if resp == nil {
					return api.ResponseServerError
				}
				return resp
			case <-ctxWithTimeout.Done():
				if errors.Is(ctxWithTimeout.Err(), context.DeadlineExceeded) {
					return api.ResponseTimeout
				}
				return api.ResponseServerError
			}
		})
	}
}
