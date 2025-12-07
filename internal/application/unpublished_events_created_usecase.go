package application

import (
	"context"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/utils"
)

type UnpublishedEventsCreatedUseCase kernel.UseCase[*kernel.NoInput, *kernel.NoOutput]

type unpublishedEventsCreatedUseCase struct {
	publisher event.Publisher
	lister    event.UnpublishedEventsLister
}

func NewUnpublishedEventsCreatedUseCase(
	publisher event.Publisher,
	lister event.UnpublishedEventsLister,
) (UnpublishedEventsCreatedUseCase, error) {
	if err := utils.CheckInterfaces(
		publisher,
		lister,
	); err != nil {
		return nil, err
	}
	return &unpublishedEventsCreatedUseCase{
		publisher: publisher,
		lister:    lister,
	}, nil
}

func (uc *unpublishedEventsCreatedUseCase) Execute(ctx context.Context, _ *kernel.NoInput) (*kernel.NoOutput, error) {
	evs, err := uc.lister.ListUnpublishedEvents(ctx)
	if err != nil {
		return nil, err
	}

	if len(evs) != 0 {
		for _, ev := range evs {
			if err := uc.publisher.Publish(ev); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
