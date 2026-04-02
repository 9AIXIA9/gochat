package gorm

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"
)

type MysqlConfig struct {
	Host                string `mapstructure:"Host"`
	Port                int    `mapstructure:"Port"`
	Username            string `mapstructure:"Username"`
	Password            string `mapstructure:"Password"`
	Database            string `mapstructure:"Database"`
	SlowThresholdMillis int    `mapstructure:"SlowThresholdMillis"`
	MaxOpenConns        int    `mapstructure:"MaxOpenConns"`
	MaxIdleConns        int    `mapstructure:"MaxIdleConns"`
	ConnMaxLifetimeSec  int    `mapstructure:"ConnMaxLifetimeSec"`
	ConnMaxIdleTimeSec  int    `mapstructure:"ConnMaxIdleTimeSec"`
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
	if c.SlowThresholdMillis < 0 || c.SlowThresholdMillis > 60_000 {
		return fmt.Errorf("%w: Mysql.SlowThresholdMillis must be in [0,60000], got %d", myErrors.ErrInvalidNumber, c.SlowThresholdMillis)
	}
	if c.MaxOpenConns < 0 {
		return fmt.Errorf("%w: Mysql.MaxOpenConns must be >= 0, got %d", myErrors.ErrInvalidNumber, c.MaxOpenConns)
	}
	if c.MaxIdleConns < 0 {
		return fmt.Errorf("%w: Mysql.MaxIdleConns must be >= 0, got %d", myErrors.ErrInvalidNumber, c.MaxIdleConns)
	}
	if c.MaxOpenConns > 0 && c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("%w: Mysql.MaxIdleConns(%d) cannot be greater than MaxOpenConns(%d)", myErrors.ErrInvalidNumber, c.MaxIdleConns, c.MaxOpenConns)
	}
	if c.ConnMaxLifetimeSec < 0 {
		return fmt.Errorf("%w: Mysql.ConnMaxLifetimeSec must be >= 0, got %d", myErrors.ErrInvalidNumber, c.ConnMaxLifetimeSec)
	}
	if c.ConnMaxIdleTimeSec < 0 {
		return fmt.Errorf("%w: Mysql.ConnMaxIdleTimeSec must be >= 0, got %d", myErrors.ErrInvalidNumber, c.ConnMaxIdleTimeSec)
	}
	return nil
}
