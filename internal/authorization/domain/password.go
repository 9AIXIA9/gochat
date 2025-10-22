package domain

import (
	"gochat/internal/shared/errors"
	"gochat/internal/shared/kernel"
)

var _ kernel.Validatable = Password("")

const (
	passwordMinLength = 6
	passwordMaxLength = 20
)

type Password string

func (p Password) Hash(encryptor kernel.HashEncryptor) (string, error) {
	return encryptor.Encrypt(string(p))
}

func (p Password) String() string {
	return string(p)
}

func (p Password) Validate() error {
	if l := len(p); l < passwordMinLength || l > passwordMaxLength {
		return errors.ErrInvalidLength
	}
	return nil
}
