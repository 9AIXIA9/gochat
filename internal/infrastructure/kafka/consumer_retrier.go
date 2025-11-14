package kafka

import (
	"context"
	"fmt"
	"gochat/internal/shared/event"
	"gochat/pkg/utils"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	maxConsumerRetries = 3
	consumerMaxBackoff = 1 * time.Second
)

type ConsumerRetrier struct {
	producer          *ProducerWithRetry
	deadLetterSaver   event.DeadLetterSaver
	publishResultChan chan kafka.Event
}

func NewConsumerRetrier(
	producer *ProducerWithRetry,
	deadLetterSaver event.DeadLetterSaver,
) *ConsumerRetrier {
	return &ConsumerRetrier{
		producer:          producer,
		deadLetterSaver:   deadLetterSaver,
		publishResultChan: make(chan kafka.Event, 512),
	}
}

func (r *ConsumerRetrier) Retry(ctx context.Context, ev event.Event, reason error, retry int) error {
	// 判断错误是否可重试
	if !r.isRetriableError() {
		if err := r.deadLetterSaver.SaveDeadLetter(ctx, ev, reason); err != nil {
			return err
		}
		return fmt.Errorf("non-retriable error: %w", reason)
	}

	if retry > maxConsumerRetries {
		zap.L().Warn("event reached max retries, sending to dead letter", zap.String("event_id", ev.ID().String()))
		if err := r.deadLetterSaver.SaveDeadLetter(ctx, ev, reason); err != nil {
			return err
		}
		return fmt.Errorf("event reached max retries: %w", reason)
	}

	utils.BackoffWait(consumerMaxBackoff, retry)

	message := getMessage(ev, retry+1)
	if err := r.producer.produceWithRetry(message, r.publishResultChan); err != nil {
		zap.L().Warn("event republish failed, sending to dead letter", zap.String("event_id", ev.ID().String()))
		if err := r.deadLetterSaver.SaveDeadLetter(ctx, ev, reason); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (r *ConsumerRetrier) isRetriableError() bool {
	return true
}

// Start background delivery handling.
func (r *ConsumerRetrier) Start() {
	go r.processSendingResponse()
}

// Close flushes and closes the producer and delivery channel.
func (r *ConsumerRetrier) Close() {
	// Producer won't write to our private delivery channel after Close/Flush
	close(r.publishResultChan)
}

func (r *ConsumerRetrier) processSendingResponse() {
	for result := range r.publishResultChan {
		switch message := result.(type) {
		case *kafka.Message:
			if message.TopicPartition.Error != nil {
				continue
			}

			ev, retry := parseMessage(message)
			zap.L().Info(
				"event retried successfully",
				zap.Int("retry_count", retry),
				zap.String("event_id", ev.ID().String()),
			)
		case kafka.Error:
			zap.L().Error("kafka producer error", zap.Error(message))
		default:
			// ignore
		}
	}
}
