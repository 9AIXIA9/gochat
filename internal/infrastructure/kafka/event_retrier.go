package kafka

import (
	"context"
	"fmt"
	"gochat/internal/shared/event"

	"go.uber.org/zap"
)

const maxRetries = 3

type EventRetrier struct {
	publisher       event.Publisher
	deadLetterSaver event.DeadLetterSaver
}

func NewEventRetrier(
	publisher event.Publisher,
	deadLetterSaver event.DeadLetterSaver,
) *EventRetrier {
	return &EventRetrier{
		publisher:       publisher,
		deadLetterSaver: deadLetterSaver,
	}
}

func (r *EventRetrier) Retry(ctx context.Context, ev event.Event, reason error) error {
	// 判断错误是否可重试
	if !r.isRetriableError() {
		if err := r.deadLetterSaver.SaveDeadLetter(ctx, ev, reason); err != nil {
			return err
		}
		return fmt.Errorf("non-retriable error: %w", reason)
	}

	retryTimes := getRetryMeta(ev.Payload())

	if retryTimes > maxRetries {
		zap.L().Warn("event reached max retries, sending to dead letter", zap.String("event_id", ev.ID().String()))
		if err := r.deadLetterSaver.SaveDeadLetter(ctx, ev, reason); err != nil {
			return err
		}
		return fmt.Errorf("event reached max retries: %w", reason)
	}

	retryTimes++

	if err := r.publisher.Publish(wrapEventWithRetryMeta(ev, retryTimes)); err != nil {
		zap.L().Warn("event republish failed, sending to dead letter", zap.String("event_id", ev.ID().String()))
		if err := r.deadLetterSaver.SaveDeadLetter(ctx, ev, reason); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (r *EventRetrier) isRetriableError() bool {
	return true
}
