package config

import (
	myErrors "gochat/internal/shared/errors"
	"time"
)

type Outbox struct {
	RunTimeout    time.Duration `mapstructure:"RunTimeout"`
	Lease         time.Duration `mapstructure:"Lease"`
	BatchSize     int           `mapstructure:"BatchSize"`
	SweepInterval time.Duration `mapstructure:"SweepInterval"`
	TriggerBuffer int           `mapstructure:"TriggerBuffer"`
}

func (c *Outbox) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.RunTimeout <= 0 {
		c.RunTimeout = time.Minute
	}
	if c.Lease <= 0 {
		c.Lease = time.Minute
	}
	if c.BatchSize <= 0 {
		c.BatchSize = 200
	}
	if c.SweepInterval <= 0 {
		c.SweepInterval = time.Second
	}
	if c.TriggerBuffer <= 0 {
		c.TriggerBuffer = 1
	}
	return nil
}
