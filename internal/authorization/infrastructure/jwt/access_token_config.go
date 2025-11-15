package jwt

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
	"time"
)

type AccessTokenConfig struct {
	Secret           string        `mapstructure:"Secret"`
	ValidityDuration time.Duration `mapstructure:"ValidityDuration"`
}

func (c *AccessTokenConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Secret == "" {
		return fmt.Errorf("%w: Secret is empty", myErrors.ErrEmptyInput)
	}
	if c.ValidityDuration <= 0 {
		return fmt.Errorf("%w: ValidityDuration must be > 0, got %s", myErrors.ErrInvalidNumber, c.ValidityDuration)
	}
	return nil
}
