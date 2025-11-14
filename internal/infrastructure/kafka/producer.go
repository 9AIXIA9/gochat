package kafka

import (
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(
	commonConfig *CommonConfig,
	producerConfig *ProducerConfig,
) (*ckafka.Producer, error) {
	kafkaConfig, err := getProducerConfig(commonConfig, producerConfig)
	if err != nil {
		return nil, err
	}

	producer, err := ckafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer failed: %w", err)
	}
	return producer, nil
}
