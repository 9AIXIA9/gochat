package kafka

import (
	"context"
	"fmt"
	"gochat/internal/infrastructure/metrics"
	myErrors "gochat/internal/shared/errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.AsyncPublisher = (*EventAsyncPublisher)(nil)

type EventAsyncPublisher struct {
	publishResultChan chan ckafka.Event
	producer          *ckafka.Producer
	onDelivered       func(event.ID) error

	wg     sync.WaitGroup
	closed atomic.Bool
	// deliveredCh decouples delivery-report handling from the kafka events loop.
	// Workers read from this channel and persist publish state concurrently to avoid
	// blocking the producer event loop (which may cause "Queue full").
	deliveredCh chan event.ID
	// number of concurrent workers that call onDelivered
	deliveredWorkers int
}

func NewEventAsyncPublisher(
	config *Config,
	onDelivered func(event.ID) error,
) (*EventAsyncPublisher, error) {
	if onDelivered == nil {
		return nil, fmt.Errorf("%w: onDelivered is nil", myErrors.ErrEmptyPointer)
	}

	producer, err := ckafka.NewProducer(getProducerConfigMap(config))
	if err != nil {
		defer func() {
			if producer == nil {
				return
			}
			producer.Close()
		}()
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	return &EventAsyncPublisher{
		// 增大结果通道，减少在高吞吐下消费者处理不及导致的阻塞
		publishResultChan: make(chan ckafka.Event, 4096),
		producer:          producer,
		onDelivered:       onDelivered,
		deliveredCh:       make(chan event.ID, 16384),
		deliveredWorkers:  16,
	}, nil
}

func (p *EventAsyncPublisher) Publish(ctx context.Context, ev event.Event) error {
	if ev == nil {
		return myErrors.ErrEmptyPointer
	}

	if p.closed.Load() {
		return myErrors.ErrHasBeenClosed
	}

	message := toMessage(ev)

	// 当本地队列已满时，进行有上下文感知的重试（指数退避），以避免把 Queue full 直接上抛到上层批量处理逻辑
	backoff := 10 * time.Millisecond
	const maxBackoff = 500 * time.Millisecond
	for {
		// respect context cancellation first
		select {
		case <-ctx.Done():
			metrics.KafkaProduce(ctx, "failed", "context_done")
			return fmt.Errorf("publish event timeout: %w", ctx.Err())
		default:
		}

		if err := p.producer.Produce(message, p.publishResultChan); err != nil {
			metrics.KafkaProduce(ctx, "failed", "produce_error")
			if strings.Contains(err.Error(), "Queue full") {
				// 记录并上报指标，然后休眠后重试
				metrics.KafkaQueueFull(ctx)
				zap.L().Warn("kafka producer queue full, retrying", zap.Duration("backoff", backoff))

				timer := time.NewTimer(backoff)
				select {
				case <-ctx.Done():
					timer.Stop()
					metrics.KafkaProduce(ctx, "failed", "context_done")
					return fmt.Errorf("publish event timeout: %w", ctx.Err())
				case <-timer.C:
				}

				if backoff < maxBackoff {
					backoff *= 2
					if backoff > maxBackoff {
						backoff = maxBackoff
					}
				}
				continue
			}
			zap.L().Error("produce message failed", zap.Error(err))
			return err
		}

		metrics.KafkaProduce(ctx, "ok", "")
		return nil
	}
}

func (p *EventAsyncPublisher) processPublishingResponse() {
	for result := range p.publishResultChan {
		switch message := result.(type) {
		case *ckafka.Message:
			if err := message.TopicPartition.Error; err != nil {
				continue
			}
			// 绝不丢弃 delivered id：当通道写满时在这里背压等待，保证最终会执行 onDelivered。
			// 这样会降低峰值吞吐，但能避免已投递消息因未标记 published 导致重复处理。
			p.deliveredCh <- message.Opaque.(event.ID)
		case ckafka.Error:
			zap.L().Error("kafka producer error", zap.Error(message))
		default:
			// ignore
		}
	}
	// publishResultChan 已关闭，通知 delivered workers 也可以结束
	close(p.deliveredCh)
}

func (p *EventAsyncPublisher) Start() {
	// start publisher event loop
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.processPublishingResponse()
	}()

	// start delivered workers
	for i := 0; i < p.deliveredWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for id := range p.deliveredCh {
				if err := p.onDelivered(id); err != nil {
					zap.L().Error("kafka message delivered callback error",
						zap.String("event_id", id.String()),
						zap.Error(err),
					)
				}
			}
		}()
	}
}

func (p *EventAsyncPublisher) Close() {
	if !p.closed.CompareAndSwap(false, true) {
		return
	}

	// 尽量冲刷剩余消息; 超时时间可配置
	p.producer.Flush(10_000)

	// 关闭通道让处理 goroutine 退出
	close(p.publishResultChan)

	// 等待 goroutine 完成
	p.wg.Wait()

	// 关闭底层 producer
	p.producer.Close()
}

func toMessage(ev event.Event) *ckafka.Message {
	topic := ev.Topic().String()

	headers := []ckafka.Header{
		{
			Key:   eventIDKey,
			Value: []byte(ev.ID()),
		},
	}

	for key, value := range ev.Headers() {
		headers = append(headers, ckafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	return &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value:     ev.Payload(),
		Timestamp: ev.OccurredAt(),
		Key:       []byte(ev.AggregateID().String()),
		Headers:   headers,
		Opaque:    ev.ID(), // 回调 识别消息
	}
}
