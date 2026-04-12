package middleware

import (
	"context"
	"fmt"
	"gochat/internal/infrastructure/kafka"
	"gochat/pkg/retry"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	retryHeaderKey = "retry"
	maxRetry       = 3
)

func NewRetryErrorMiddleware(
	producer *ckafka.Producer,
) kafka.ErrorMiddleware {
	return func(next kafka.ErrorHandler) kafka.ErrorHandler {
		return kafka.ErrorHandlerFunc(func(ctx context.Context, err error, message *ckafka.Message) {
			//进行重试
			var retryCount int
			found := false
			for i, header := range message.Headers {
				if header.Key == retryHeaderKey {
					found = true
					if _, err := fmt.Sscanf(string(header.Value), "%d", &retryCount); err != nil {
						retryCount = 0
					}
					message.Headers[i].Value = []byte(fmt.Sprintf("%d", retryCount+1))
					break
				}
			}
			if !found {
				// initialize retry header
				message.Headers = append(message.Headers, ckafka.Header{Key: retryHeaderKey, Value: []byte("1")})
				retryCount = 1
			}

			if retryCount >= maxRetry {
				next.Handle(ctx, err, message)
				return
			}

			//复制消息并增加重试次数
			newMessage := &ckafka.Message{
				TopicPartition: message.TopicPartition,
				Key:            message.Key,
				Value:          message.Value,
				Headers:        append(message.Headers, ckafka.Header{Key: retryHeaderKey, Value: []byte(fmt.Sprintf("%d", retryCount+1))}),
			}

			//发送消息到Kafka
			err = producer.Produce(newMessage, nil)
			if err != nil {
				retry.Wait(ctx, retryCount, 100*time.Millisecond)
				//如果发送失败，调用下一个错误处理器
				next.Handle(ctx, err, message)
				return
			}
			zap.L().Info(
				"kafka message retried",
				zap.Int("retry_count", retryCount),
				zap.String("topic", *message.TopicPartition.Topic),
				zap.Binary("key", message.Key),
			)
		})
	}
}
