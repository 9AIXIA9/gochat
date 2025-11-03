package config

import (
	"fmt"
	"gochat/internal/authorization/infrastructure/bcrypt"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	"gochat/internal/delivery/http/middleware"
	"gochat/internal/infrastructure/canal"
	"gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/zap"
	"gochat/internal/notification/infrastructure/gomail"
	myErrors "gochat/internal/shared/errors"
)

// App 应用程序配置结构
type App struct {
	Name        string `mapstructure:"Name"`
	Host        string `mapstructure:"Host"`
	Port        int    `mapstructure:"Port"`
	MachineNode int64  `mapstructure:"MachineNode"`

	Cookie       *Cookie                     `mapstructure:"Cookie"`
	CORS         *middleware.CORSConfig      `mapstructure:"CORS"`
	AccessToken  *jwt.AccessTokenConfig      `mapstructure:"AccessToken"`
	RefreshToken *crypto.RefreshTokenConfig  `mapstructure:"RefreshToken"`
	Hasher       *bcrypt.HasherConfig        `mapstructure:"Hasher"`
	Mysql        *gorm.MysqlConfig           `mapstructure:"Mysql"`
	Redis        *redis.Config               `mapstructure:"Redis"`
	Kafka        *kafka.Config               `mapstructure:"Kafka"`
	Logger       *zap.LoggerConfig           `mapstructure:"Logger"`
	RateLimit    *middleware.RateLimitConfig `mapstructure:"RateLimit"`
	BinlogReader *canal.BinlogReaderConfig   `mapstructure:"BinlogReader"`
	Email        *gomail.EmailNotifierConfig `mapstructure:"Email"`
}

func (c *App) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}
	if c.Name == "" {
		return fmt.Errorf("App.Name: %w", myErrors.ErrEmptyInput)
	}
	if c.Host == "" {
		return fmt.Errorf("App.Host: %w", myErrors.ErrEmptyInput)
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("App.Port: %w: must be in 1..65535, got %d", myErrors.ErrInvalidNumber, c.Port)
	}
	// Allow 0 as valid MachineNode (snowflake supports node 0..1023)
	if c.MachineNode < 0 || c.MachineNode > 1023 {
		return fmt.Errorf("App.MachineNode: %w: must be in 0..1023, got %d", myErrors.ErrInvalidNumber, c.MachineNode)
	}

	if err := c.Cookie.Validate(); err != nil {
		return fmt.Errorf("App.Cookie: %w", err)
	}
	if err := c.CORS.Validate(); err != nil {
		return fmt.Errorf("App.CORS: %w", err)
	}
	if err := c.AccessToken.Validate(); err != nil {
		return fmt.Errorf("App.AccessToken: %w", err)
	}
	if err := c.RefreshToken.Validate(); err != nil {
		return fmt.Errorf("App.RefreshToken: %w", err)
	}
	if err := c.Mysql.Validate(); err != nil {
		return fmt.Errorf("App.Mysql: %w", err)
	}
	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("App.Redis: %w", err)
	}
	if err := c.Kafka.Validate(); err != nil {
		return fmt.Errorf("App.Kafka: %w", err)
	}
	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("App.Logger: %w", err)
	}
	if err := c.RateLimit.Validate(); err != nil {
		return fmt.Errorf("App.RateLimit: %w", err)
	}
	if err := c.Email.Validate(); err != nil {
		return fmt.Errorf("App.Email: %w", err)
	}
	if err := c.BinlogReader.Validate(); err != nil {
		return fmt.Errorf("App.BinlogReader: %w", err)
	}

	return nil
}
