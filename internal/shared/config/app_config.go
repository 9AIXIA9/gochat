package config

import (
	"fmt"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	myErrors "gochat/internal/shared/errors"
	"gochat/internal/shared/http/middlewares"
	"gochat/internal/shared/infrastructure/bcrypt"
	"gochat/internal/shared/infrastructure/gorm"
	"gochat/internal/shared/infrastructure/redis"
	"gochat/internal/shared/infrastructure/zap"
)

// App 应用程序配置结构
type App struct {
	Name        string `mapstructure:"Name"`
	Host        string `mapstructure:"Host"`
	Port        int    `mapstructure:"Port"`
	MachineNode int64  `mapstructure:"MachineNode"`

	Cert         *Cert                        `mapstructure:"Cert"`
	Cookie       *Cookie                      `mapstructure:"Cookie"`
	CORS         *middlewares.CORSConfig      `mapstructure:"CORS"`
	AccessToken  *jwt.AccessTokenConfig       `mapstructure:"AccessToken"`
	RefreshToken *crypto.RefreshTokenConfig   `mapstructure:"RefreshToken"`
	Hasher       *bcrypt.HasherConfig         `mapstructure:"Hasher"`
	Mysql        *gorm.MysqlConfig            `mapstructure:"Mysql"`
	Redis        *redis.Config                `mapstructure:"Redis"`
	Logger       *zap.LoggerConfig            `mapstructure:"Logger"`
	RateLimit    *middlewares.RateLimitConfig `mapstructure:"RateLimit"`
}

func (c *App) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Name == "" {
		return fmt.Errorf("%w: Name is empty", myErrors.ErrEmptyInput)
	}
	if c.Host == "" {
		return fmt.Errorf("%w: Host is empty", myErrors.ErrEmptyInput)
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("%w: Port must be in 1..65535, got %d", myErrors.ErrInvalidNumber, c.Port)
	}
	// Allow 0 as valid MachineNode (snowflake supports node 0..1023)
	if c.MachineNode < 0 || c.MachineNode > 1023 {
		return fmt.Errorf("%w: MachineNode must be in 0..1023, got %d", myErrors.ErrInvalidNumber, c.MachineNode)
	}

	// Validate sub-configs when provided
	if c.Cert != nil {
		if err := c.Cert.Validate(); err != nil {
			return err
		}
	}
	if c.Cookie != nil {
		if err := c.Cookie.Validate(); err != nil {
			return err
		}
	}
	if c.CORS != nil {
		if err := c.CORS.Validate(); err != nil {
			return err
		}
	}
	if c.AccessToken != nil {
		if err := c.AccessToken.Validate(); err != nil {
			return err
		}
	}
	if c.RefreshToken != nil {
		if err := c.RefreshToken.Validate(); err != nil {
			return err
		}
	}
	if c.Mysql != nil {
		if err := c.Mysql.Validate(); err != nil {
			return err
		}
	}
	if c.Redis != nil {
		if err := c.Redis.Validate(); err != nil {
			return err
		}
	}
	if c.Logger != nil {
		if err := c.Logger.Validate(); err != nil {
			return err
		}
	}
	return nil
}
