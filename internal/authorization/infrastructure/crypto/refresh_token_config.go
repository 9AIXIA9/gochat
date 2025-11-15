package crypto

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type RefreshTokenConfig struct {
	Length int `mapstructure:"Length"`
}

func (c *RefreshTokenConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Length <= 0 {
		return fmt.Errorf("%w: RefreshTokenConfig.Length must be > 0, got %d", myErrors.ErrInvalidNumber, c.Length)
	}
	return nil
}
