package application

import (
	"context"
	"errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const processDuration = 1 * time.Minute

type UnpublishedEventsCreatedUseCase kernel.UseCase[*kernel.NoInput, *kernel.NoOutput]

type unpublishedEventsCreatedUseCase struct {
	publisher event.Publisher
	lister    event.UnpublishedEventsLister
	running   atomic.Bool
}

func NewUnpublishedEventsCreatedUseCase(
	publisher event.Publisher,
	lister event.UnpublishedEventsLister,
) (UnpublishedEventsCreatedUseCase, error) {
	if err := validate.NotNil(
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
	if !uc.running.CompareAndSwap(false, true) {
		zap.L().Debug("Skip outbox processing because a previous execution is still running")
		return nil, nil
	}
	defer uc.running.Store(false)

	evs, err := uc.lister.ListUnpublishedEvents(ctx, processDuration)
	if err != nil {
		return nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, processDuration)
	defer cancel()

	if len(evs) != 0 {
		for _, ev := range evs {
			if err := uc.publisher.Publish(ctxWithTimeout, ev); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					zap.L().Info("Deadline exceeded while publishing events, stopping further attempts")
					return nil, nil
				}
				return nil, err
			}
		}
	}
	return nil, nil
}
