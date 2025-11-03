package canal

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type BinlogReaderConfig struct {
	Addr     string `mapstructure:"Addr"`
	User     string `mapstructure:"User"`
	Password string `mapstructure:"Password"`
	TableDB  string `mapstructure:"TableDB"`
}

func (c *BinlogReaderConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Addr == "" {
		return fmt.Errorf("%w: BinlogReaderConfig.Addr is empty", myErrors.ErrEmptyInput)
	}
	if c.User == "" {
		return fmt.Errorf("%w: BinlogReaderConfig.User is empty", myErrors.ErrEmptyInput)
	}
	if c.Password == "" {
		return fmt.Errorf("%w: BinlogReaderConfig.Password is empty", myErrors.ErrEmptyInput)
	}
	if c.TableDB == "" {
		return fmt.Errorf("%w: BinlogReaderConfig.TableDB is empty", myErrors.ErrEmptyInput)
	}
	return nil
}
