package kafka

import (
	"context"
	"fmt"
	myErrors "gochat/internal/shared/errors"
	"gochat/pkg/utils"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

const (
	pollTimeout = 500 * time.Millisecond
)

type ErrorHandler interface {
	Handle(ctx context.Context, err error, message *ckafka.Message)
}
type ErrorHandlerFunc func(ctx context.Context, err error, message *ckafka.Message)

type Consumer struct {
	consumer     *ckafka.Consumer
	router       *Router
	errorHandler ErrorHandler

	ctx    context.Context
	cancel context.CancelFunc

	running bool
}

func NewConsumer(config *Config, router *Router) (*Consumer, error) {
	consumer, err := ckafka.NewConsumer(getConsumerConfigMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer failed: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		consumer:     consumer,
		router:       router,
		errorHandler: nil,
		ctx:          ctx,
		cancel:       cancel,
		running:      false,
	}, nil
}

func (c *Consumer) Start() error {
	if c.running {
		return nil
	}
	topics := c.router.Topics()
	if len(topics) == 0 {
		return fmt.Errorf("kafka consumer router topics is empty %w", myErrors.ErrEmptyInput)
	}

	if err := c.consumer.SubscribeTopics(topics, nil); err != nil {
		return fmt.Errorf("subscribe topics failed: %w", err)
	}

	c.running = true
	zap.L().Info("kafka consumer started", zap.Strings("topics", topics), zap.Bool("auto_commit", true))

	utils.GoSafe(c.processMessage)

	return nil
}

func (c *Consumer) processMessage() {
	for c.running {
		// Handle ctx cancellation
		select {
		case <-c.ctx.Done():
			c.running = false
			return
		default:
		}

		ev := c.consumer.Poll(int(pollTimeout.Milliseconds()))
		if ev == nil {
			continue
		}

		switch m := ev.(type) {
		case *ckafka.Message:
			// Route the message
			err := c.router.Route(c.ctx, m)
			if err != nil {
				// Log / handle the processing error
				if c.errorHandler != nil {
					c.errorHandler.Handle(c.ctx, err, m)
				} else {
					zap.L().Error(
						"kafka consumer handle message failed",
						zap.String("topic", *m.TopicPartition.Topic),
						zap.Int32("partition", m.TopicPartition.Partition),
						zap.Int64("offset", int64(m.TopicPartition.Offset)),
						zap.Error(err),
					)
				}
			}
		case ckafka.Error:
			// Consumer-level error
			zap.L().Error("kafka consumer error", zap.Error(m))
		default:
			// ignore other events
		}
	}
}

func (c *Consumer) SetErrorHandler(h ErrorHandler) {
	c.errorHandler = h
}

func (c *Consumer) Close() {
	c.running = false
	c.cancel()
	if c.consumer != nil {
		_ = c.consumer.Close()
	}
}
