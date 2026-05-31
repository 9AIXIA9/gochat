package kafka

import (
	"context"
	"fmt"
	"gochat/internal/shared/command"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

var _ command.SyncPublisher = (*CommandSyncPublisher)(nil)

type CommandSyncPublisher struct {
	producer *ckafka.Producer
}

func NewCommandSyncPublisher(config *Config) (*CommandSyncPublisher, error) {
	producer, err := ckafka.NewProducer(getProducerConfigMap(config))
	if err != nil {
		defer func() {
			if producer == nil {
				return
			}
			producer.Close()
		}()
		return nil, fmt.Errorf("create kafka sync producer failed: %w", err)
	}

	return &CommandSyncPublisher{
		producer: producer,
	}, nil
}

func (p *CommandSyncPublisher) Publish(ctx context.Context, com command.Command) error {
	if com == nil {
		return fmt.Errorf("command is nil")
	}

	message := commandToMessage(com)

	// 创建一个只给当前这次投递使用的一次性接收 channel
	deliveryChan := make(chan ckafka.Event, 1)
	defer close(deliveryChan)

	// 调用异步 Produce 方法，但我们把自己的 deliveryChan 传进去做回调监听
	err := p.producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	// 阻塞当前协程，等待消息真正的发送结果或者 ctx 超时
	select {
	case <-ctx.Done():
		return fmt.Errorf("sync publish context done: %w", ctx.Err())

	case e := <-deliveryChan:
		m, ok := e.(*ckafka.Message)
		if !ok {
			return fmt.Errorf("unexpected command type received from delivery channel")
		}
		// m.TopicPartition.Error 如果非 nil，说明远端明确拒绝了或者超时失败等
		if err := m.TopicPartition.Error; err != nil {
			return fmt.Errorf("kafka broker rejected message: %w", err)
		}
		// 如果过了这一行，说明消息确实安全抵达 Broker ！
		return nil
	}
}

func (p *CommandSyncPublisher) Close() {
	if p.producer != nil {
		p.producer.Flush(5000)
		p.producer.Close()
	}
}
