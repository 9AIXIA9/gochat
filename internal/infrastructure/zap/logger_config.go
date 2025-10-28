package zap

import (
	"fmt"
	"strings"

	myErrors "gochat/internal/shared/errors"
)

type LoggerConfig struct {
	Mode       string `mapstructure:"Mode"`
	Level      string `mapstructure:"Level"`
	Filename   string `mapstructure:"Filename"`
	MaxSize    int    `mapstructure:"MaxSize"`
	MaxAge     int    `mapstructure:"MaxAge"`
	MaxBackups int    `mapstructure:"MaxBackups"`
}

func (c *LoggerConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if strings.TrimSpace(c.Mode) == "" {
		return fmt.Errorf("%w: Log.Mode is empty", myErrors.ErrEmptyInput)
	}
	if strings.TrimSpace(c.Level) == "" {
		return fmt.Errorf("%w: Log.Level is empty", myErrors.ErrEmptyInput)
	}
	if strings.TrimSpace(c.Filename) == "" {
		return fmt.Errorf("%w: Log.Filename is empty", myErrors.ErrEmptyInput)
	}
	if c.MaxSize <= 0 {
		return fmt.Errorf("%w: Log.MaxSize must be > 0, got %d", myErrors.ErrInvalidNumber, c.MaxSize)
	}
	if c.MaxAge < 0 {
		return fmt.Errorf("%w: Log.MaxAge must be >= 0, got %d", myErrors.ErrInvalidNumber, c.MaxAge)
	}
	if c.MaxBackups < 0 {
		return fmt.Errorf("%w: Log.MaxBackups must be >= 0, got %d", myErrors.ErrInvalidNumber, c.MaxBackups)
	}
	return nil
}
