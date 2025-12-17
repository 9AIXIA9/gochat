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
	if l := len(p); l < passwordMinLength {
		return errors.WrapBusiness(errors.ErrInvalidLength, "password too short: min %d", passwordMinLength)
	} else if l > passwordMaxLength {
		return errors.WrapBusiness(errors.ErrInvalidLength, "password too long: max %d", passwordMaxLength)
	}
	return nil
}

func (p Password) Encrypt(encryptor Encryptor) (PasswordEncrypted, error) {
	if len(p) == 0 {
		return "", errors.WrapBusiness(errors.ErrEmptyInput, "password is empty")
	}
	encrypted, err := encryptor.Encrypt(p.String())
	if err != nil {
		return "", err
	}
	return PasswordEncrypted(encrypted), nil
}

func (e PasswordEncrypted) Compare(
	password Password,
	comparator Comparator,
) error {
	if len(e) == 0 {
		return nil
	}

	if err := comparator.Compare(e.String(), password.String()); err != nil {
		return ErrInvalidPassword
	}
	return nil
}
