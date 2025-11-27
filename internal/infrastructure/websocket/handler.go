package websocket

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context, data []byte) ([]byte, error)
}

type HandlerFunc func(ctx context.Context, data []byte) ([]byte, error)

func (f HandlerFunc) Handle(ctx context.Context, data []byte) ([]byte, error) {
	return f(ctx, data)
}

func chainHandlers(mainHandler Handler, middlewares []Handler) Handler {
	if len(middlewares) == 0 {
		return mainHandler
	}

	return HandlerFunc(func(ctx context.Context, msg []byte) ([]byte, error) {
		var err error
		currentMsg := msg

		// 依次执行中间件
		for _, middleware := range middlewares {
			currentMsg, err = middleware.Handle(ctx, currentMsg)
			if err != nil {
				return nil, err
			}
		}

		// 最后执行主处理程序
		return mainHandler.Handle(ctx, currentMsg)
	})
}
