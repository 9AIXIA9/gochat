package event

import "context"

type Repository interface {
	UnpublishedSaver
	UnpublishedLister
	PublishedMarker
	DeadLetterSaver
}

type UnpublishedSaver interface {
	Saves(ctx context.Context, events []Event) error
}

type UnpublishedLister interface {
	UnpublishedList(ctx context.Context) ([]Event, error)
}

type PublishedMarker interface {
	MarkPublished(ctx context.Context, ID ID) error
}

type DeadLetterSaver interface {
	SaveDeadLetter(ctx context.Context, event Event, reason error) error
}
