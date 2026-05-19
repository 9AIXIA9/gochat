package zap

import (
	"fmt"
	"strings"

	myErrors "gochat/internal/shared/errors"
)

type LoggerConfig struct {
	Level string `mapstructure:"Level"`
}

func (c *LoggerConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if strings.TrimSpace(c.Level) == "" {
		return fmt.Errorf("%w: Log.Level is empty", myErrors.ErrEmptyInput)
	}
	return nil
}
