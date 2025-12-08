//go:generate mockgen -source=repository.go -destination=./mocks/mock_repository.go -package=mocks
package event

import "context"

type Repository interface {
	UnpublishedEventsCreator
	UnpublishedEventsLister
	PublishedAndNotProcessingMarker
	NotProcessingMarker
	DeadLetterCreator
}

type UnpublishedEventsCreator interface {
	CreateUnpublishedEvents(ctx context.Context, evs []Event) error
}

type UnpublishedEventsLister interface {
	ListUnpublishedEvents(ctx context.Context) ([]Event, error)
}

type PublishedAndNotProcessingMarker interface {
	MarkAsPublishedAndNotProcessing(ctx context.Context, ID ID) error
}

type NotProcessingMarker interface {
	MarkAsNotProcessing(ctx context.Context, ID ID) error
}

type DeadLetterCreator interface {
	CreateDeadLetter(ctx context.Context, e Event, reason error) error
}
