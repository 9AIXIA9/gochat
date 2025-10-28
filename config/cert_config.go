package config

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type Cert struct {
	HTTPSKeyFile  string `mapstructure:"HTTPSKeyFile"`
	HTTPSCertFile string `mapstructure:"HTTPSCertFile"`
}

func (c *Cert) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.HTTPSKeyFile == "" {
		return fmt.Errorf("%w: Cert.HTTPSKeyFile is empty", myErrors.ErrEmptyInput)
	}
	if c.HTTPSCertFile == "" {
		return fmt.Errorf("%w: Cert.HTTPSCertFile is empty", myErrors.ErrEmptyInput)
	}
	return nil
}
