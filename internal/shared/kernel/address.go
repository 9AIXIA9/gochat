package kernel

import (
	"gochat/internal/shared/errors"
)

const maxAddressLen = 100

type Address string

func (a Address) String() string {
	return string(a)
}

// Validate 验证地址格式
func (a Address) Validate() error {
	if l := len(a); l == 0 {
		return errors.ErrEmptyInput
	} else if l > maxAddressLen {
		return errors.ErrInvalidLength
	}
	return nil
}
