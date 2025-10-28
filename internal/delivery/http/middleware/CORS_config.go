package middleware

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
	"time"
)

type CORSConfig struct {
	Origins          []string      `mapstructure:"Origins"`
	AllowMethods     []string      `mapstructure:"AllowMethods"`
	AllowHeaders     []string      `mapstructure:"AllowHeaders"`
	ExposeHeaders    []string      `mapstructure:"ExposeHeaders"`
	AllowCredentials bool          `mapstructure:"AllowCredentials"`
	MaxAge           time.Duration `mapstructure:"MaxAge"`
}

func (c *CORSConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.MaxAge < 0 {
		return fmt.Errorf("%w: NewCORSMiddleware.MaxAge must be >= 0, got %s", myErrors.ErrInvalidNumber, c.MaxAge)
	}
	return nil
}
