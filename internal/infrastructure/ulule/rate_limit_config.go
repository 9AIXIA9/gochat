package ulule

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
	"time"
)

type Config struct {
	Period time.Duration `mapstructure:"Period"`
	Limit  int64         `mapstructure:"Limit"`
}

func (c *Config) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Period < 0 {
		return fmt.Errorf("%w: Config.Period must be >= 0, got %s", myErrors.ErrInvalidNumber, c.Period)
	}
	if c.Limit < 0 {
		return fmt.Errorf("%w: Config.Limit must be >= 0, got %d", myErrors.ErrInvalidNumber, c.Limit)
	}
	return nil
}
