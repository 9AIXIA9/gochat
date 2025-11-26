package application

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
)

type UnpublishedEventsCreatedUseCase kernel.UseCase[*kernel.NoInput, *kernel.NoOutput]

type unpublishedEventsCreatedUseCase struct {
	publisher event.Publisher
	lister    event.UnpublishedEventsLister
}

func NewUnpublishedEventsCreatedUseCase(
	publisher event.Publisher,
	lister event.UnpublishedEventsLister,
) UnpublishedEventsCreatedUseCase {
	return &unpublishedEventsCreatedUseCase{
		publisher: publisher,
		lister:    lister,
	}
}

func (uc *unpublishedEventsCreatedUseCase) Execute(ctx context.Context, _ *kernel.NoInput) (*kernel.NoOutput, error) {
	evs, err := uc.lister.ListUnpublishedEvents(ctx)
	if err != nil {
		return nil, err
	}

	for _, ev := range evs {
		if err := uc.publisher.Publish(ev); err != nil {
			return nil, err
		}

	}
	return nil, nil
}
