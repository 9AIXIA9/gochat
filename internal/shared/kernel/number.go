package kernel

import (
	"gochat/pkg/utils"
)

var _ Validatable = Number("")

type Number string

func (n Number) Validate() error {
	return utils.ValidateNumber(string(n))
}
