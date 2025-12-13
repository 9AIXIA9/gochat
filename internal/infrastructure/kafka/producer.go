package kafka

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(conf *Config) (*ckafka.Producer, error) {
	producer, err := ckafka.NewProducer(getProducerConfigMap(conf))
	if err != nil {
		defer func() {
			if producer == nil {
				return
			}
			producer.Close()
		}()
		return nil, err
	}
	return producer, nil
}
