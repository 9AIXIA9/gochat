package command

import (
	"context"
)

var _ Handler = HandlerFunc(nil)

type Handler interface {
	Handle(ctx context.Context, c Command) error
}

type HandlerFunc func(ctx context.Context, c Command) error

func (h HandlerFunc) Handle(ctx context.Context, c Command) error {
	return h(ctx, c)
}
