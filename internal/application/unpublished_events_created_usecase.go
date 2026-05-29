package application

import (
	"context"
	"errors"
	"gochat/internal/shared/event"
	"gochat/internal/shared/kernel"
	"gochat/pkg/validate"
	"time"

	"go.uber.org/zap"
)

const (
	defaultProcessDuration = 1 * time.Minute
	defaultEventsLimit     = 200
)

type UnpublishedEventsCreatedUseCase kernel.UseCase[*kernel.NoInput, *kernel.NoOutput]

type unpublishedEventsCreatedUseCase struct {
	publisher  event.AsyncPublisher
	lister     event.UnpublishedEventsLister
	lease      time.Duration
	limit      int
	runTimeout time.Duration
}

func NewUnpublishedEventsCreatedUseCase(
	publisher event.AsyncPublisher,
	lister event.UnpublishedEventsLister,
) (UnpublishedEventsCreatedUseCase, error) {
	return NewUnpublishedEventsCreatedUseCaseWithOptions(publisher, lister, 0, 0, 0)
}

func NewUnpublishedEventsCreatedUseCaseWithOptions(
	publisher event.AsyncPublisher,
	lister event.UnpublishedEventsLister,
	lease time.Duration,
	limit int,
	runTimeout time.Duration,
) (UnpublishedEventsCreatedUseCase, error) {
	if err := validate.NotNil(
		publisher,
		lister,
	); err != nil {
		return nil, err
	}
	if lease <= 0 {
		lease = defaultProcessDuration
	}
	if limit <= 0 {
		limit = defaultEventsLimit
	}
	if runTimeout <= 0 {
		runTimeout = defaultProcessDuration
	}
	return &unpublishedEventsCreatedUseCase{
		publisher:  publisher,
		lister:     lister,
		lease:      lease,
		limit:      limit,
		runTimeout: runTimeout,
	}, nil
}

func (uc *unpublishedEventsCreatedUseCase) Execute(ctx context.Context, _ *kernel.NoInput) (*kernel.NoOutput, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, uc.runTimeout)
	defer cancel()

	for {
		evs, err := uc.lister.ListUnpublishedEvents(ctxWithTimeout, uc.lease, uc.limit)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				zap.L().Info("Deadline exceeded while listing unpublished events, stopping further attempts")
				return nil, nil
			}
			return nil, err
		}

		if len(evs) == 0 {
			return nil, nil
		}

		for _, ev := range evs {
			if err := uc.publisher.Publish(ctxWithTimeout, ev); err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					zap.L().Info("Deadline exceeded while publishing events, stopping further attempts")
					return nil, nil
				}
				return nil, err
			}
		}

		if len(evs) < uc.limit {
			return nil, nil
		}
	}
}
