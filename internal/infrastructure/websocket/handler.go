package websocket

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context, data []byte) ([]byte, error)
}

// Middleware wraps a next Handler and returns a new Handler.
type Middleware func(next Handler) Handler

type HandlerFunc func(ctx context.Context, data []byte) ([]byte, error)

func (f HandlerFunc) Handle(ctx context.Context, data []byte) ([]byte, error) {
	return f(ctx, data)
}

// chainHandlers wraps the mainHandler with provided middlewares (outermost first).
func chainHandlers(mainHandler Handler, middlewares []Middleware) Handler {
	if len(middlewares) == 0 {
		return mainHandler
	}

	// Apply in reverse so the first middleware becomes the outermost wrapper.
	wrapped := mainHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
