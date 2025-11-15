package kernel

import (
	"gochat/internal/shared/errors"
)

var _ Validatable = Number("")

type Number string

func (n Number) String() string {
	return string(n)
}

func (n Number) Validate() error {
	for _, digit := range n {
		if digit < '0' || digit > '9' {
			return errors.ErrInvalidNumber
		}
	}
	return nil
}
