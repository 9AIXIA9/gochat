package domain

import "gochat/internal/shared/kernel"

var _ kernel.Validatable = UserNumber("")

type UserNumber kernel.Number

func (n UserNumber) String() string {
	return string(n)
}

func (n UserNumber) Validate() error {
	return kernel.Number(n).Validate()
}
