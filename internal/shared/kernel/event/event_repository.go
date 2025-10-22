package event

import "context"

type Saver interface {
	Save(ctx context.Context, events []Event) error
}

type PublisherRepository[SpecificEvent Event] interface {
	UnpublishedEventLister[SpecificEvent]
	PublishStatusUpdater
}

type UnpublishedEventLister[SpecificEvent Event] interface {
	ListUnpublishedEvent(ctx context.Context) ([]SpecificEvent, error)
}

type PublishStatusUpdater interface {
	MarkAsPublished(ctx context.Context, IDs []ID) error
}
