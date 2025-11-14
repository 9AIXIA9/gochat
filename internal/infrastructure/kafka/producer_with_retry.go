package kafka

import (
	"errors"
	"gochat/pkg/utils"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	maxProducerRetries = 3
	maxProducerBackoff = 1 * time.Second
)

type ProducerWithRetry struct {
	producer *ckafka.Producer
}

func NewProducerWithRetry(producer *ckafka.Producer) *ProducerWithRetry {
	return &ProducerWithRetry{
		producer: producer,
	}
}

func (p *ProducerWithRetry) produceWithRetry(
	msg *ckafka.Message,
	resultChan chan ckafka.Event,
) error {
	var lastErr error
	for attempt := 0; attempt <= maxProducerRetries; attempt++ {
		lastErr = p.producer.Produce(msg, resultChan)
		if lastErr == nil {
			return nil
		}
		var kerr ckafka.Error
		if errors.As(lastErr, &kerr) && kerr.Code() == ckafka.ErrQueueFull {
			// back off and let delivery goroutine drain
			utils.BackoffWait(maxProducerBackoff, attempt)
			continue
		}
		// non-retryable error
		break
	}
	return lastErr
}

func (p *ProducerWithRetry) Close() {
	_ = p.producer.Flush(int((3 * time.Second).Milliseconds()))
	p.producer.Close()
}
