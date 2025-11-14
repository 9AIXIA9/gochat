package kafka

import (
	"context"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// EnsureTopics creates topics if they do not already exist. It is safe to call repeatedly.
// numPartitions <= 0 defaults to 1; replicationFactor <= 0 defaults to 1.
func EnsureTopics(ctx context.Context, config *Config, topics []string, numPartitions int, replicationFactor int) error {
	if config == nil {
		return fmt.Errorf("kafka ensure topics: config is nil")
	}
	if len(topics) == 0 {
		return nil
	}

	if numPartitions <= 0 {
		numPartitions = 1
	}
	if replicationFactor <= 0 {
		replicationFactor = 1
	}

	admin, err := ckafka.NewAdminClient(convertToMap(config))
	if err != nil {
		return fmt.Errorf("create kafka admin client failed: %w", err)
	}
	defer admin.Close()

	// de-dup topics
	seen := make(map[string]struct{}, len(topics))
	specs := make([]ckafka.TopicSpecification, 0, len(topics))
	for _, t := range topics {
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		specs = append(specs, ckafka.TopicSpecification{
			Topic:             t,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		})
	}
	if len(specs) == 0 {
		return nil
	}

	results, err := admin.CreateTopics(ctx, specs)
	if err != nil {
		return fmt.Errorf("create topics request failed: %w", err)
	}
	for _, r := range results {
		if r.Error.Code() == ckafka.ErrNoError || r.Error.Code() == ckafka.ErrTopicAlreadyExists {
			continue
		}
		return fmt.Errorf("create topic %s failed: %v", r.Topic, r.Error)
	}
	return nil
}
