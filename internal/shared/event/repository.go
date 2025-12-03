//go:generate mockgen -source=repository.go -destination=./mocks/mock_repository.go -package=mocks
package event

import "context"

type Repository interface {
	UnpublishedEventCreator
	UnpublishedEventsCreator
	UnpublishedEventsLister
	PublishedMarker
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

type PublishedMarker interface {
	MarkAsPublished(ctx context.Context, ID ID) error
}

type DeadLetterCreator interface {
	CreateDeadLetter(ctx context.Context, e Event, reason error) error
}
