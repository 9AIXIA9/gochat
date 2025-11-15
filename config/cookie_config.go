package config

import (
	"fmt"
	"strings"

	myErrors "gochat/internal/shared/errors"
)

type Cookie struct {
	Domain   string `mapstructure:"Domain"`
	Secure   bool   `mapstructure:"Secure"`
	HttpOnly bool   `mapstructure:"HttpOnly"`
	Path     string `mapstructure:"Path"`
}

func (c *Cookie) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if strings.TrimSpace(c.Path) == "" {
		return fmt.Errorf("%w: Cookie.Path is empty", myErrors.ErrEmptyInput)
	}
	// No further strict checks; Domain may be empty, Secure/HttpOnly booleans are fine.
	return nil
}
