package kafka

import (
	kafkautil "gochat/internal/infrastructure/kafka"
)

func NewSubscriber(
	commonConfig *kafkautil.CommonConfig,
	consumerConfig *kafkautil.ConsumerConfig,
) (*kafkautil.EventSubscriber, error) {
	kafkaSubscriber, err := kafkautil.NewEventSubscriber(commonConfig, consumerConfig)
	if err != nil {
		return nil, err
	}

	return kafkaSubscriber, nil
}
