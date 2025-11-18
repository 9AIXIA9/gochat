package event

import "context"

type Repository interface {
	UnpublishedEventSaver
	UnpublishedEventsSaver
	UnpublishedLister
	PublishedMarker
	DeadLetterSaver
}

type UnpublishedEventSaver interface {
	Save(ctx context.Context, event Event) error
}

type UnpublishedEventsSaver interface {
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
