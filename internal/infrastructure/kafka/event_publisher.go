package kafka

import (
	"context"
	"fmt"
	myErrors "gochat/internal/shared/errors"
	"sync"
	"sync/atomic"

	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/shared/event"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var _ event.Publisher = (*EventPublisher)(nil)

type EventPublisher struct {
	publishResultChan chan ckafka.Event
	producer          *ckafka.Producer
	publishedMarker   event.PublishedMarker //TODO 感觉不该让publisher来使用 而是上层调用者来使用
	metrics           *prometheus.Metrics

	wg     sync.WaitGroup
	closed atomic.Bool
}

func NewEventPublisher(
	config *Config,
	publisher event.PublishedMarker,
	metrics *prometheus.Metrics,
) (*EventPublisher, error) {
	producer, err := ckafka.NewProducer(getProducerConfigMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	return &EventPublisher{
		publishResultChan: make(chan ckafka.Event, 512),
		producer:          producer,
		publishedMarker:   publisher,
		metrics:           metrics,
	}, nil
}

func (p *EventPublisher) publish(event event.Event) error {
	if event == nil {
		return myErrors.ErrEmptyPointer
	}

	if p.closed.Load() {
		return myErrors.ErrHasBeenClosed
	}

	message := p.getMessage(event)
	if err := p.producer.Produce(message, p.publishResultChan); err != nil {
		if p.metrics != nil {
			p.metrics.KafkaProduced.WithLabelValues(*message.TopicPartition.Topic, "error").Inc()
		}
		zap.L().Error("produce message failed", zap.Error(err))
		return err
	}
	return nil
}

func (p *EventPublisher) Publish(events []event.Event) error {
	if len(events) == 0 {
		return nil
	}
	var firstErr error
	for _, e := range events {
		if err := p.publish(e); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (p *EventPublisher) processPublishingResponse() {
	for result := range p.publishResultChan {
		switch message := result.(type) {
		case *ckafka.Message:
			if message.TopicPartition.Error != nil {
				if p.metrics != nil {
					p.metrics.KafkaProduced.WithLabelValues(*message.TopicPartition.Topic, "error").Inc()
				}
				continue
			}

			id := message.Opaque.(event.ID)
			if err := p.publishedMarker.MarkAsPublished(context.Background(), id); err != nil {
				zap.L().Error(
					"mark event published failed",
					zap.Error(err),
					zap.String("event_id", id.String()),
				)
			} else if p.metrics != nil {
				p.metrics.KafkaProduced.WithLabelValues(*message.TopicPartition.Topic, "success").Inc()
			}
		case ckafka.Error:
			zap.L().Error("kafka producer error", zap.Error(message))
		default:
			// ignore
		}
	}
}

func (p *EventPublisher) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.processPublishingResponse()
	}()
}

func (p *EventPublisher) Close() {
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

func (p *EventPublisher) getMessage(ev event.Event) *ckafka.Message {
	topic := ev.Topic().String()
	m := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          ev.Payload(),
		Timestamp:      ev.OccurredAt(),
		Key:            []byte(ev.AggregateID().String()),
		Headers: []ckafka.Header{
			{
				Key:   "event_id",
				Value: []byte(ev.ID()),
			},
		},
		Opaque: ev.ID(), // 回调 识别消息
	}
	return m
}
