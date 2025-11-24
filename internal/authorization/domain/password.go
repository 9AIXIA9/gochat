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

type (
	Password          string
	PasswordEncrypted string
)

func (p Password) String() string {
	return string(p)
}

func (e PasswordEncrypted) String() string {
	return string(e)
}

func (p Password) Validate() error {
	if l := len(p); l < passwordMinLength || l > passwordMaxLength {
		return errors.ErrInvalidLength
	}
	return nil
}

func (p Password) Encrypt(encryptor Encryptor) (PasswordEncrypted, error) {
	encrypted, err := encryptor.Encrypt(p.String())
	if err != nil {
		return "", err
	}
	return PasswordEncrypted(encrypted), nil
}
