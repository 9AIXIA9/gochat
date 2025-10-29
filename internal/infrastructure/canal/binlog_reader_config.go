package canal

import (
	"fmt"
	myErrors "gochat/internal/shared/errors"

	"github.com/go-mysql-org/go-mysql/canal"
)

type BinlogReaderConfig struct {
	Addr          string `mapstructure:"Addr"`
	User          string `mapstructure:"User"`
	Password      string `mapstructure:"Password"`
	TableDB       string `mapstructure:"TableDB"`
	ExecutionPath string `mapstructure:"ExecutionPath"`
}

func (c *BinlogReaderConfig) ToCanal() *canal.Config {
	config := canal.NewDefaultConfig()
	config.Addr = c.Addr
	config.User = c.User
	config.Password = c.Password
	config.Dump.TableDB = c.TableDB
	config.Dump.ExecutionPath = c.ExecutionPath
	return config
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
	//c.ExecutionPath == ""  是允许的 -> 不使用 mysqldump 工具
	return nil
}
