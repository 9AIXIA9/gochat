package gomail

import (
	"fmt"
	"gochat/internal/shared/errors"
)

type EmailNotifierConfig struct {
	Host     string `mapstructure:"Host"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
	Port     int    `mapstructure:"Port"`
}

func (c *EmailNotifierConfig) Validate() error {
	if c == nil {
		return errors.ErrEmptyPointer
	}
	if c.Host == "" {
		return fmt.Errorf("%w: EmailNotifierConfig.Host is empty", errors.ErrEmptyInput)
	}
	if c.Username == "" {
		return fmt.Errorf("%w: EmailNotifierConfig.Username is empty", errors.ErrEmptyInput)
	}
	if c.Password == "" {
		return fmt.Errorf("%w: EmailNotifierConfig.Password is empty", errors.ErrEmptyInput)
	}
	if c.Port == 0 {
		return fmt.Errorf("%w: EmailNotifierConfig.Port is empty", errors.ErrEmptyInput)
	}
	return nil
}
