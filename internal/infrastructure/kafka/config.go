package kafka

import (
	"fmt"
	"strings"
)

type Config struct {
	Common   *CommonConfig   `mapstructure:"Common"`
	Consumer *ConsumerConfig `mapstructure:"Consumer"`
	Producer *ProducerConfig `mapstructure:"Producer"`
}

type CommonConfig struct {
	BootstrapServers string `mapstructure:"BootstrapServers"`
	ClientID         string `mapstructure:"ClientID"`
	SecurityProtocol string `mapstructure:"SecurityProtocol"`
	SASLMechanism    string `mapstructure:"SASLMechanism"`
	SASLUsername     string `mapstructure:"SASLUsername"`
	SASLPassword     string `mapstructure:"SASLPassword"`
}

type ConsumerConfig struct {
	GroupID              string `mapstructure:"GroupID"`
	AutoOffsetReset      string `mapstructure:"AutoOffsetReset"`
	EnableAutoCommit     *bool  `mapstructure:"EnableAutoCommit"`
	AutoCommitIntervalMs int    `mapstructure:"AutoCommitIntervalMs"`
	SessionTimeoutMs     int    `mapstructure:"SessionTimeoutMs"`
	MaxPollIntervalMs    int    `mapstructure:"MaxPollIntervalMs"`
}

type ProducerConfig struct {
	Acks                  string `mapstructure:"Acks"`
	EnableIdempotence     bool   `mapstructure:"EnableIdempotence"`
	MessageTimeoutMs      int    `mapstructure:"MessageTimeoutMs"`
	AllowAutoCreateTopics bool   `mapstructure:"AllowAutoCreateTopics"`
	MaxInFlight           int    `mapstructure:"MaxInFlight"`
	MaxRetries            int    `mapstructure:"MaxRetries"`
	BackoffMs             int    `mapstructure:"BackoffMs"`
}

// ---- Validate ----

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("kafka config: nil pointer")
	}
	if err := c.Common.Validate(); err != nil {
		return fmt.Errorf("common: %w", err)
	}
	if err := c.Consumer.Validate(); err != nil {
		return fmt.Errorf("consumer: %w", err)
	}
	if err := c.Producer.Validate(); err != nil {
		return fmt.Errorf("producer: %w", err)
	}
	return nil
}

func (c *CommonConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("common config: nil pointer")
	}
	if strings.TrimSpace(c.BootstrapServers) == "" {
		return fmt.Errorf("bootstrap servers is required")
	}
	sp := strings.ToLower(strings.TrimSpace(c.SecurityProtocol))
	isSASL := sp == "sasl_plaintext" || sp == "sasl_ssl"
	if isSASL {
		if strings.TrimSpace(c.SASLMechanism) == "" {
			return fmt.Errorf("SASL mechanism is required when using SASL security protocol")
		}
		if strings.TrimSpace(c.SASLUsername) == "" || strings.TrimSpace(c.SASLPassword) == "" {
			return fmt.Errorf("SASL username/password are required when using SASL security protocol")
		}
	}
	return nil
}

func (c *ConsumerConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("consumer config: nil pointer")
	}
	if strings.TrimSpace(c.GroupID) == "" {
		return fmt.Errorf("group.id is required")
	}
	aor := strings.ToLower(strings.TrimSpace(c.AutoOffsetReset))
	if aor != "earliest" && aor != "latest" {
		return fmt.Errorf("auto.offset.reset must be earliest or latest")
	}
	if c.EnableAutoCommit == nil {
		return fmt.Errorf("enable.auto.commit is required")
	}
	// 其它数值边界由 tag 约束；此处不强制默认值
	return nil
}

func (p *ProducerConfig) Validate() error {
	if p == nil {
		return fmt.Errorf("producer config: nil pointer")
	}
	acks := strings.TrimSpace(p.Acks)
	if acks != "all" && acks != "-1" && acks != "0" && acks != "1" {
		return fmt.Errorf("acks must be one of all|-1|0|1")
	}
	// 幂等性要求：MaxInFlight<=5（librdkafka 推荐）
	if p.EnableIdempotence && p.MaxInFlight > 5 {
		return fmt.Errorf("max.in.flight.requests.per.connection must be <= 5 when enable.idempotence=true")
	}
	return nil
}
