package kafka

import (
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewConsumer(
	commonConfig *CommonConfig,
	consumerConfig *ConsumerConfig,
) (*ckafka.Consumer, error) {
	cCfg, err := getConsumerConfig(commonConfig, consumerConfig)
	if err != nil {
		return nil, err
	}

	consumer, err := ckafka.NewConsumer(cCfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka consumer failed: %w", err)
	}

	return consumer, nil
}
