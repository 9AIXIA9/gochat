package event

import "context"

type Repository interface {
	UnpublishedEventCreator
	UnpublishedEventsCreator
	UnpublishedEventsLister
	Publisher
	DeadLetterCreator
}

type UnpublishedEventCreator interface {
	CreateUnpublishedEvent(ctx context.Context, e Event) error
}

type UnpublishedEventsCreator interface {
	CreateUnpublishedEvents(ctx context.Context, evs []Event) error
}

type UnpublishedEventsLister interface {
	ListUnpublishedEvents(ctx context.Context) ([]Event, error)
}

type Publisher interface {
	Publish(ctx context.Context, ID ID) error
}

type DeadLetterCreator interface {
	CreateDeadLetter(ctx context.Context, e Event, reason error) error
}
