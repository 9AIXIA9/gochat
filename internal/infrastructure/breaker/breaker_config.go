package breaker

import "time"

type Config struct {
	Name           string        `mapstructure:"Name"`
	MaxRequests    uint32        `mapstructure:"MaxRequests"`
	Interval       time.Duration `mapstructure:"Interval"`
	BucketPeriod   time.Duration `mapstructure:"BucketPeriod"`
	Timeout        time.Duration `mapstructure:"Timeout"`
	ErrorRate      float64       `mapstructure:"ErrorRate"`
	ErrorThreshold uint32        `mapstructure:"ErrorThreshold"`
}

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}
	if c.Name == "" {
		c.Name = "default_breaker"
	}
	if c.MaxRequests == 0 {
		c.MaxRequests = 5
	}
	if c.Interval == 0 {
		c.Interval = time.Second * 60
	}
	if c.BucketPeriod == 0 {
		c.BucketPeriod = time.Second * 10
	}
	if c.Timeout == 0 {
		c.Timeout = time.Second * 10
	}
	if c.ErrorRate == 0 {
		c.ErrorRate = 0.6
	}
	if c.ErrorThreshold == 0 {
		c.ErrorThreshold = 10
	}
	return nil
}
