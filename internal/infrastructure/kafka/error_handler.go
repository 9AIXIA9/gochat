package kafka

import (
	"context"
	"fmt"
	"gochat/internal/shared/event"
	"gochat/pkg/utils"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	retryHeaderKey = "retry"
	maxRetry       = 3
)

type ErrorHandlerWithDeadLetterAndRetry struct {
	isRetriableError func(error) bool
	producer         *ckafka.Producer
	creator          event.DeadLetterCreator
}

func NewErrorHandlerWithDeadLetterAndRetry(
	config *Config,
	creator event.DeadLetterCreator,
	isRetriableError func(error) bool,
) (ErrorHandler, error) {
	if config == nil {
		return nil, fmt.Errorf("kafka error handler: config is nil")
	}
	if creator == nil {
		return nil, fmt.Errorf("kafka error handler: dead letter creator is nil")
	}
	if isRetriableError == nil {
		return nil, fmt.Errorf("kafka error handler: retry judge is nil")
	}
	producer, err := ckafka.NewProducer(getProducerConfigMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka producer for error handler failed: %w", err)
	}
	return &ErrorHandlerWithDeadLetterAndRetry{
		isRetriableError: isRetriableError,
		producer:         producer,
		creator:          creator,
	}, nil
}

func (h *ErrorHandlerWithDeadLetterAndRetry) Handle(ctx context.Context, err error, message *ckafka.Message) {
	var retry int
	found := false
	for i, header := range message.Headers {
		if header.Key == retryHeaderKey {
			found = true
			if _, err := fmt.Sscanf(string(header.Value), "%d", &retry); err != nil {
				retry = 0
			}
			message.Headers[i].Value = []byte(fmt.Sprintf("%d", retry+1))
			break
		}
	}
	if !found {
		// initialize retry header
		message.Headers = append(message.Headers, ckafka.Header{Key: retryHeaderKey, Value: []byte("1")})
		retry = 1
	}

	if retry <= maxRetry && h.isRetriableError(err) {
		h.reproduce(ctx, err, message, retry)
	} else {
		if err := h.creator.CreateDeadLetter(ctx, toEvent(message), err); err != nil {
			zap.L().Error("create dead letter failed", zap.Error(err))
		}
	}
}

func (h *ErrorHandlerWithDeadLetterAndRetry) reproduce(ctx context.Context, err error, message *ckafka.Message, retry int) {
	if h.producer == nil {
		zap.L().Error("retry reproduce failed: producer is nil", zap.Error(err))
		return
	}

	utils.BackoffWait(ctx, retry, 500*time.Millisecond)

	// Re-produce to the same topic/partition (PartitionAny lets broker decide)
	if err := h.producer.Produce(message, nil); err != nil {
		zap.L().Error(
			"kafka reproduce message failed",
			zap.String("topic", *message.TopicPartition.Topic),
			zap.Int32("partition", message.TopicPartition.Partition),
			zap.Int64("offset", int64(message.TopicPartition.Offset)),
			zap.Error(err),
		)
		return
	}
	zap.L().Info(
		"kafka message scheduled for retry",
		zap.String("topic", *message.TopicPartition.Topic),
		zap.Int32("partition", message.TopicPartition.Partition),
		zap.Int64("offset", int64(message.TopicPartition.Offset)),
		zap.Int("retry", retry),
	)
}
