package kafka

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

// Config wraps producer-related settings for Kafka
// Only commonly used fields are exposed; extend as needed.
type Config struct {
	BootstrapServers      string `mapstructure:"BootstrapServers"`
	ClientID              string `mapstructure:"ClientID"`
	Acks                  string `mapstructure:"Acks"`
	SecurityProtocol      string `mapstructure:"SecurityProtocol"`
	SASLMechanism         string `mapstructure:"SASLMechanism"`
	SASLUsername          string `mapstructure:"SASLUsername"`
	SASLPassword          string `mapstructure:"SASLPassword"`
	EnableIdempotence     bool   `mapstructure:"EnableIdempotence"`
	MessageTimeoutMs      int    `mapstructure:"MessageTimeoutMs"`
	AllowAutoCreateTopics bool   `mapstructure:"AllowAutoCreateTopics"`
	MaxInFlight           int    `mapstructure:"MaxInFlight"` // in-flight requests per connection
	MaxRetries            int    `mapstructure:"MaxRetries"`  // client-side retry on produce failure (queue full etc.)
	BackoffMs             int    `mapstructure:"BackoffMs"`   // backoff between local retries
}

func (c *Config) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.BootstrapServers == "" {
		return fmt.Errorf("Kafka.BootstrapServers: %w", myErrors.ErrEmptyInput)
	}
	// optional fields are fine
	return nil
}

// toKafkaConfig converts to confluent-kafka-go ConfigMap
func (c *Config) toKafkaConfig() (*ckafka.ConfigMap, error) {
	cfg := &ckafka.ConfigMap{
		"bootstrap.servers":        c.BootstrapServers,
		"allow.auto.create.topics": c.AllowAutoCreateTopics,
	}
	if c.ClientID != "" {
		if err := cfg.SetKey("client.id", c.ClientID); err != nil {
			return nil, fmt.Errorf("Kafka.ClientID: %w", err)
		}
	}
	if c.Acks != "" {
		if err := cfg.SetKey("acks", c.Acks); err != nil {
			return nil, fmt.Errorf("Kafka.Acks: %w", err)
		}
	}
	if c.EnableIdempotence {
		if err := cfg.SetKey("enable.idempotence", true); err != nil {
			return nil, fmt.Errorf("Kafka.EnableIdempotence: %w", err)
		}
	}
	if c.MessageTimeoutMs > 0 {
		if err := cfg.SetKey("message.timeout.ms", c.MessageTimeoutMs); err != nil {
			return nil, fmt.Errorf("Kafka.MessageTimeoutMs: %w", err)
		}
	}
	if c.MaxInFlight > 0 {
		if err := cfg.SetKey("max.in.flight.requests.per.connection", c.MaxInFlight); err != nil {
			return nil, fmt.Errorf("Kafka.MaxInFlight: %w", err)
		}
	}
	// Security (optional)
	if c.SecurityProtocol != "" {
		if err := cfg.SetKey("security.protocol", c.SecurityProtocol); err != nil {
			return nil, fmt.Errorf("Kafka.SecurityProtocol: %w", err)
		}
	}
	if c.SASLMechanism != "" {
		if err := cfg.SetKey("sasl.mechanism", c.SASLMechanism); err != nil {
			return nil, fmt.Errorf("Kafka.SASLMechanism: %w", err)
		}
	}
	if c.SASLUsername != "" {
		if err := cfg.SetKey("sasl.username", c.SASLUsername); err != nil {
			return nil, fmt.Errorf("Kafka.SASLUsername: %w", err)
		}
	}
	if c.SASLPassword != "" {
		if err := cfg.SetKey("sasl.password", c.SASLPassword); err != nil {
			return nil, fmt.Errorf("Kafka.SASLPassword: %w", err)
		}
	}
	return cfg, nil
}
