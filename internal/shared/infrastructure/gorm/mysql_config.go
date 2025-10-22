package gorm

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type MysqlConfig struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
	Database string `mapstructure:"Database"`
}

func (c *MysqlConfig) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Host == "" {
		return fmt.Errorf("%w: Mysql.Host is empty", myErrors.ErrEmptyInput)
	}
	if c.Username == "" {
		return fmt.Errorf("%w: Mysql.Username is empty", myErrors.ErrEmptyInput)
	}
	if c.Database == "" {
		return fmt.Errorf("%w: Mysql.Database is empty", myErrors.ErrEmptyInput)
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("%w: Mysql.Port must be in 1..65535, got %d", myErrors.ErrInvalidNumber, c.Port)
	}
	return nil
}
