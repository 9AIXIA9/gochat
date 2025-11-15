package domain

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

var _ kernel.Validatable = Password("")

const (
	passwordMaxLength = 20
)

type Password string

func (p Password) String() string {
	return string(p)
}

func (p Password) Validate() error {
	if l := len(p); l > passwordMaxLength {
		return errors.ErrInvalidLength
	}
	return nil
}
