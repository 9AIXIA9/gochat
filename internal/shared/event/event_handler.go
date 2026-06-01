package event

import "context"

var _ Handler = HandlerFunc(nil)

type Handler interface {
	Handle(ctx context.Context, e Event) error
}

type HandlerFunc func(ctx context.Context, e Event) error

func (h HandlerFunc) Handle(ctx context.Context, e Event) error {
	return h(ctx, e)
}
