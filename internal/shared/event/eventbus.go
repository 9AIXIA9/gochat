package event

import "context"

type Publisher interface {
	Publish(event Event) error
	Publishes(events []Event) error
}

type Subscriber interface {
	Subscribe(topic Topic, handler Handler)
}

type Handler interface {
	Handle(ctx context.Context, e Event) error
}

type HandlerFunc func(ctx context.Context, e Event) error

func (h HandlerFunc) Handle(ctx context.Context, e Event) error {
	return h(ctx, e)
}
