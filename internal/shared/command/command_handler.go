package command

import "context"

type Handler interface {
	Handle(ctx context.Context, c Command) error
}

type HandlerFunc func(ctx context.Context, c Command) error

func (h HandlerFunc) Handle(ctx context.Context, c Command) error {
	return h(ctx, c)
}
