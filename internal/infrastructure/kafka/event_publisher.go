package kafka

import (
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
	metrics           *prometheus.Metrics
	onDelivered       func(event.ID) error

	wg     sync.WaitGroup
	closed atomic.Bool
}

func NewEventPublisher(
	config *Config,
	metrics *prometheus.Metrics,
	onDelivered func(event.ID) error,
) (*EventPublisher, error) {
	if onDelivered == nil {
		return nil, fmt.Errorf("%w: onDelivered is nil", myErrors.ErrEmptyPointer)
	}

	producer, err := ckafka.NewProducer(getProducerConfigMap(config))
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}

	return &EventPublisher{
		publishResultChan: make(chan ckafka.Event, 512),
		producer:          producer,
		metrics:           metrics,
		onDelivered:       onDelivered,
	}, nil
}

func (p *EventPublisher) Publish(event event.Event) error {
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

			if err := p.onDelivered(message.Opaque.(event.ID)); err != nil {
				zap.L().Error(
					"kafka message delivered callback error",
					zap.String("event_id", message.Opaque.(event.ID).String()),
					zap.Error(err),
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
