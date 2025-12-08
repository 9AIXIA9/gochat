package handler

import (
	"context"
	"gochat/internal/infrastructure/kafka"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func NewLoggerErrorHandler() kafka.ErrorHandlerFunc {
	return func(ctx context.Context, err error, message *ckafka.Message) {
		zap.L().Error("kafka message processing error",
			zap.String("topic", *message.TopicPartition.Topic),
			zap.Int32("partition", message.TopicPartition.Partition),
			zap.Int64("offset", int64(message.TopicPartition.Offset)),
			zap.ByteString("key", message.Key),
			zap.ByteString("value", message.Value),
			zap.Error(ctx.Err()),
		)
	}
}
