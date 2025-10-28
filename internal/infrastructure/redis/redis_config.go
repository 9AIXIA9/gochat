package redis

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type Config struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Password string `mapstructure:"Password"`
	Database int    `mapstructure:"Database"`
}

func (c *Config) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Host == "" {
		return fmt.Errorf("%w: Redis.Host is empty", myErrors.ErrEmptyInput)
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("%w: Redis.Port must be in 1..65535, got %d", myErrors.ErrInvalidNumber, c.Port)
	}
	if c.Database < 0 {
		return fmt.Errorf("%w: Redis.Database must be >= 0, got %d", myErrors.ErrInvalidNumber, c.Database)
	}
	return nil
}
