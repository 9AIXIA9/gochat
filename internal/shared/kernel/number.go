package kernel

import "gochat/pkg/utils"

var _ Validatable = Number("")

type Number string

func (n Number) String() string {
	return string(n)
}

func (n Number) Validate() error {
	return utils.ValidateNumber(n.String())
}
