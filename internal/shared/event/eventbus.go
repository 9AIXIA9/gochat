package event

import "context"

type Publisher interface {
	Publish(events []Event) error
}

type Subscriber interface {
	Subscribe(topic Topic, handler Handler)
}

type Handler func(ctx context.Context, e Event) error
