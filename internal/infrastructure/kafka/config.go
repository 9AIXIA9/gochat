package kafka

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	BootstrapServers string `mapstructure:"BootstrapServers"`
	GroupID          string `mapstructure:"GroupID"`
	Acks             string `mapstructure:"Acks"`
	AutoOffsetReset  string `mapstructure:"AutoOffsetReset"`
	EnableAutoCommit *bool  `mapstructure:"EnableAutoCommit"`
}

func (c *Config) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.BootstrapServers == "" {
		return fmt.Errorf("%w: Config.BootstrapServers is empty", myErrors.ErrEmptyInput)
	}
	if c.GroupID == "" {
		return fmt.Errorf("%w: Config.GroupID is empty", myErrors.ErrEmptyInput)
	}
	if c.Acks == "" {
		return fmt.Errorf("%w: Config.Acks is empty", myErrors.ErrEmptyInput)
	}
	if c.AutoOffsetReset == "" {
		return fmt.Errorf("%w: Config.AutoOffsetReset is empty", myErrors.ErrEmptyInput)
	}
	if c.EnableAutoCommit == nil {
		return fmt.Errorf("%w: Config.EnableAutoCommit is nil", myErrors.ErrEmptyInput)
	}
	return nil
}

func convertToMap(c *Config) *ckafka.ConfigMap {
	return &ckafka.ConfigMap{
		"bootstrap.servers":  c.BootstrapServers,
		"group.id":           c.GroupID,
		"acks":               c.Acks,
		"auto.offset.reset":  c.AutoOffsetReset,
		"enable.auto.commit": *c.EnableAutoCommit,
	}
}
