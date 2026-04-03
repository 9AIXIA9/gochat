package kafka

import (
	"context"
	"errors"
	"fmt"
	"gochat/internal/infrastructure/metrics"
	myErrors "gochat/internal/shared/errors"
	"gochat/pkg/concurrency"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sony/gobreaker/v2"
	"go.uber.org/zap"
)

const (
	pollTimeout = 500 * time.Millisecond
)

type Consumer struct {
	consumer     *ckafka.Consumer
	router       *Router
	errorHandler ErrorHandler

	ctx    context.Context
	cancel context.CancelFunc

	running bool

	// pause duration when circuit breaker is open before resuming the partition
	breakerPause time.Duration
}

func NewConsumer(config *Config, router *Router) (*Consumer, error) {
	consumer, err := ckafka.NewConsumer(getConsumerConfigMap(config))
	if err != nil {
		defer func() {
			if consumer == nil {
				return
			}
			if err := consumer.Close(); err != nil {
				zap.L().Warn("kafka consumer close failed after creation error", zap.Error(err))
			}
		}()
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
		breakerPause: 0,
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

	concurrency.GoSafe(c.processMessage)

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
			topic := "unknown"
			if m.TopicPartition.Topic != nil {
				topic = *m.TopicPartition.Topic
			}
			start := time.Now()
			// Route the message
			err := c.router.Route(c.ctx, m)
			if err != nil {
				metrics.KafkaConsume(c.ctx, "handle_failed", topic)
				metrics.KafkaHandleDuration(c.ctx, topic, "handle_failed", time.Since(start).Seconds())
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

				// If breaker is open/overloaded, pause this partition to avoid advancing the position
				if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
					c.pausePartition(m.TopicPartition)
				}
				// Do NOT commit on error
			}
			if err == nil {
				metrics.KafkaConsume(c.ctx, "ok", topic)
				metrics.KafkaHandleDuration(c.ctx, topic, "ok", time.Since(start).Seconds())
				// Commit only on successful handling
				if _, cerr := c.consumer.CommitMessage(m); cerr != nil {
					metrics.KafkaCommitFail(c.ctx, topic)
					zap.L().Warn("kafka manual commit failed", zap.Error(cerr))
				}
			}
		case ckafka.Error:
			metrics.KafkaConsumerError(c.ctx, m.Code().String())
			// Consumer-level error
			zap.L().Error("kafka consumer error", zap.Error(m))
		default:
			// ignore other events
		}
	}
}

func (c *Consumer) SetErrorHandler(h ErrorHandler, middlewares ...ErrorMiddleware) {
	c.errorHandler = chainErrorHandlers(h, middlewares)
}

func (c *Consumer) Close() {
	c.running = false
	c.cancel()
	if c.consumer != nil {
		_ = c.consumer.Close()
	}
}

// SetBreakerPause configures how long to pause partitions when the circuit is open.
func (c *Consumer) SetBreakerPause(d time.Duration) {
	c.breakerPause = d
}

// pausePartition pauses consumption for the given partition and schedules a resume after breakerPause.
func (c *Consumer) pausePartition(tp ckafka.TopicPartition) {
	if c.breakerPause <= 0 {
		return
	}
	if err := c.consumer.Pause([]ckafka.TopicPartition{tp}); err != nil {
		zap.L().Warn("kafka pause partition failed", zap.Error(err))
		return
	}
	// schedule resume
	pauseFor := c.breakerPause
	concurrency.GoSafe(func() {
		timer := time.NewTimer(pauseFor)
		defer timer.Stop()
		select {
		case <-timer.C:
			if err := c.consumer.Resume([]ckafka.TopicPartition{tp}); err != nil {
				zap.L().Warn("kafka resume partition failed", zap.Error(err))
			}
		case <-c.ctx.Done():
			return
		}
	})
}

func (c *Consumer) Ping(ctx context.Context) error {
	if c.consumer == nil {
		return errors.New("consumer is nil")
	}

	// 设置超时
	timeoutMs := 1000
	if deadline, ok := ctx.Deadline(); ok {
		if ms := int(time.Until(deadline).Milliseconds()); ms > 0 {
			timeoutMs = ms
		}
	}

	// 获取元数据检查连接
	if _, err := c.consumer.GetMetadata(nil, false, timeoutMs); err != nil {
		return fmt.Errorf("consumer metadata check failed: %w", err)
	}
	return nil
}
