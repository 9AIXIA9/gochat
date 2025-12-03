package kernel

import "gochat/internal/shared/errors"

var _ Validatable = Number("")
var _ Validatable = UserNumber("")

type Number string

func (n Number) String() string {
	return string(n)
}

func (n Number) Validate() error {
	if len(n) == 0 {
		return errors.ErrEmptyInput
	}
	for _, digit := range n {
		if digit < '0' || digit > '9' {
			return errors.ErrInvalidNumber
		}
	}
	return nil
}

type RoomNumber Number

func (n RoomNumber) String() string {
	return string(n)
}

func (n RoomNumber) Validate() error {
	return Number(n).Validate()
}

type UserNumber Number

func (n UserNumber) String() string {
	return string(n)
}

func (n UserNumber) Validate() error {
	return Number(n).Validate()
}
