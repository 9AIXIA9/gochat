package kafka

import (
	"fmt"
	"strings"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func getProducerConfig(commonConfig *CommonConfig, producerConfig *ProducerConfig) (*ckafka.ConfigMap, error) {
	if commonConfig == nil || producerConfig == nil {
		return nil, fmt.Errorf("nil config")
	}
	if err := commonConfig.Validate(); err != nil {
		return nil, err
	}
	if err := producerConfig.Validate(); err != nil {
		return nil, err
	}

	cm := &ckafka.ConfigMap{}
	if err := applyCommon(cm, commonConfig); err != nil {
		return nil, err
	}

	// 必要: acks
	if err := cm.SetKey("acks", strings.TrimSpace(producerConfig.Acks)); err != nil {
		return nil, err
	}
	// 可选: 幂等性
	if producerConfig.EnableIdempotence {
		if err := cm.SetKey("enable.idempotence", true); err != nil {
			return nil, err
		}
	}
	// 可选: 消息超时
	if producerConfig.MessageTimeoutMs > 0 {
		if err := cm.SetKey("message.timeout.ms", producerConfig.MessageTimeoutMs); err != nil {
			return nil, err
		}
	}
	// 可选: 允许自动建 topic（仅在显式为 true 时设置）
	if producerConfig.AllowAutoCreateTopics {
		if err := cm.SetKey("allow.auto.create.topics", true); err != nil {
			return nil, err
		}
	}
	// 可选: MaxInFlight
	if producerConfig.MaxInFlight > 0 {
		if err := cm.SetKey("max.in.flight.requests.per.connection", producerConfig.MaxInFlight); err != nil {
			return nil, err
		}
	}
	// 可选: 重试与退避
	if producerConfig.MaxRetries > 0 {
		if err := cm.SetKey("retries", producerConfig.MaxRetries); err != nil {
			return nil, err
		}
	}
	if producerConfig.BackoffMs > 0 {
		if err := cm.SetKey("retry.backoff.ms", producerConfig.BackoffMs); err != nil {
			return nil, err
		}
	}

	return cm, nil
}

func getConsumerConfig(common *CommonConfig, consumer *ConsumerConfig) (*ckafka.ConfigMap, error) {
	if common == nil || consumer == nil {
		return nil, fmt.Errorf("nil config")
	}
	if err := common.Validate(); err != nil {
		return nil, err
	}
	if err := consumer.Validate(); err != nil {
		return nil, err
	}

	cm := &ckafka.ConfigMap{}
	if err := applyCommon(cm, common); err != nil {
		return nil, err
	}

	// 必要
	if err := cm.SetKey("group.id", strings.TrimSpace(consumer.GroupID)); err != nil {
		return nil, err
	}
	if err := cm.SetKey("auto.offset.reset", strings.ToLower(strings.TrimSpace(consumer.AutoOffsetReset))); err != nil {
		return nil, err
	}
	// required 字段：始终下发
	if err := cm.SetKey("enable.auto.commit", *consumer.EnableAutoCommit); err != nil {
		return nil, err
	}

	// 可选
	if consumer.AutoCommitIntervalMs > 0 {
		if err := cm.SetKey("auto.commit.interval.ms", consumer.AutoCommitIntervalMs); err != nil {
			return nil, err
		}
	}
	if consumer.SessionTimeoutMs > 0 {
		if err := cm.SetKey("session.timeout.ms", consumer.SessionTimeoutMs); err != nil {
			return nil, err
		}
	}
	if consumer.MaxPollIntervalMs > 0 {
		if err := cm.SetKey("max.poll.interval.ms", consumer.MaxPollIntervalMs); err != nil {
			return nil, err
		}
	}

	return cm, nil
}

// ---- helpers ----

func applyCommon(cm *ckafka.ConfigMap, c *CommonConfig) error {
	// 必要: bootstrap.servers
	bs := strings.TrimSpace(c.BootstrapServers)
	if err := cm.SetKey("bootstrap.servers", bs); err != nil {
		return err
	}
	// 可选: client.id
	if id := strings.TrimSpace(c.ClientID); id != "" {
		if err := cm.SetKey("client.id", id); err != nil {
			return err
		}
	}
	// 可选: 安全协议
	if sp := strings.TrimSpace(c.SecurityProtocol); sp != "" {
		if err := cm.SetKey("security.protocol", normalizeSecurityProtocol(sp)); err != nil {
			return err
		}
	}
	// 可选: SASL
	if mech := strings.TrimSpace(c.SASLMechanism); mech != "" {
		if err := cm.SetKey("sasl.mechanism", mech); err != nil {
			return err
		}
	}
	if u := strings.TrimSpace(c.SASLUsername); u != "" {
		if err := cm.SetKey("sasl.username", u); err != nil {
			return err
		}
	}
	if p := strings.TrimSpace(c.SASLPassword); p != "" {
		if err := cm.SetKey("sasl.password", p); err != nil {
			return err
		}
	}
	return nil
}

func normalizeSecurityProtocol(sp string) string {
	switch strings.ToLower(strings.TrimSpace(sp)) {
	case "plaintext":
		return "PLAINTEXT"
	case "ssl":
		return "SSL"
	case "sasl_plaintext":
		return "SASL_PLAINTEXT"
	case "sasl_ssl":
		return "SASL_SSL"
	default:
		return sp
	}
}
