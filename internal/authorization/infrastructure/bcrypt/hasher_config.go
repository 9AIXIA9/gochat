package bcrypt

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type HasherConfig struct {
	Cost int `mapstructure:"Cost"`
}

func (c *HasherConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Cost != 0 && (c.Cost < 4 || c.Cost > 31) {
		return fmt.Errorf("%w: HasherConfig.Cost must be in 4..31 (or 0 for default), got %d", myErrors.ErrInvalidNumber, c.Cost)
	}
	return nil
}
