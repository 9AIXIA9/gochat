package config

import (
	"fmt"
	"gochat/internal/authorization/infrastructure/crypto"
	"gochat/internal/authorization/infrastructure/jwt"
	"gochat/internal/delivery/http/middleware"
	"gochat/internal/infrastructure/bcrypt"
	"gochat/internal/infrastructure/breaker"
	"gochat/internal/infrastructure/canal"
	"gochat/internal/infrastructure/gorm"
	"gochat/internal/infrastructure/kafka"
	"gochat/internal/infrastructure/otel"
	"gochat/internal/infrastructure/redis"
	"gochat/internal/infrastructure/ulule"
	"gochat/internal/infrastructure/zap"
	myErrors "gochat/internal/shared/errors"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
	defaultEnv     = "debug"
)

// App 应用程序配置结构
type App struct {
	Name                       string        `mapstructure:"Name"`
	Env                        string        `mapstructure:"Env"`
	Timeout                    time.Duration `mapstructure:"Timeout"`
	Host                       string        `mapstructure:"Host"`
	Port                       int           `mapstructure:"Port"`
	MachineNode                int64         `mapstructure:"MachineNode"`
	DisableSessionStartedEvent bool          `mapstructure:"DisableSessionStartedEvent"`
	PProf                      *PProf        `mapstructure:"PProf"`

	Cookie             *Cookie                    `mapstructure:"Cookie"`
	CORS               *middleware.CORSConfig     `mapstructure:"CORS"`
	AccessToken        *jwt.AccessTokenConfig     `mapstructure:"AccessToken"`
	RefreshToken       *crypto.RefreshTokenConfig `mapstructure:"RefreshToken"`
	Hasher             *bcrypt.HasherConfig       `mapstructure:"Hasher"`
	Mysql              *gorm.MysqlConfig          `mapstructure:"Mysql"`
	Redis              *redis.Config              `mapstructure:"Redis"`
	Kafka              *kafka.Config              `mapstructure:"Kafka"`
	Logger             *zap.LoggerConfig          `mapstructure:"Logger"`
	HTTPRateLimit      *ulule.Config              `mapstructure:"HTTPRateLimit"`
	WebsocketRateLimit *ulule.Config              `mapstructure:"WebsocketRateLimit"`
	KafkaRateLimit     *ulule.Config              `mapstructure:"KafkaRateLimit"`
	BinlogReader       *canal.BinlogReaderConfig  `mapstructure:"BinlogReader"`
	Outbox             *Outbox                    `mapstructure:"Outbox"`
	Breaker            *breaker.Config            `mapstructure:"Breaker"`
	OTEL               *otel.Config               `mapstructure:"OTEL"`
}

type PProf struct {
	Enabled bool   `mapstructure:"Enabled"`
	Host    string `mapstructure:"Host"`
	Port    int    `mapstructure:"Port"`
}

func (c *App) Validate() error {
	if c == nil {
		return myErrors.ErrEmptyPointer
	}

	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}

	if len(c.Env) == 0 {
		c.Env = defaultEnv
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

	if c.PProf != nil && c.PProf.Enabled {
		if c.PProf.Host == "" {
			return fmt.Errorf("App.PProf.Host: %w", myErrors.ErrEmptyInput)
		}
		if c.PProf.Port <= 0 || c.PProf.Port > 65535 {
			return fmt.Errorf("App.PProf.Port: %w: must be in 1..65535, got %d", myErrors.ErrInvalidNumber, c.PProf.Port)
		}
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
	if err := c.HTTPRateLimit.Validate(); err != nil {
		return fmt.Errorf("App.HTTPRatelimit: %w", err)
	}
	if err := c.WebsocketRateLimit.Validate(); err != nil {
		return fmt.Errorf("App.WebsocketRateLimit: %w", err)
	}
	if err := c.KafkaRateLimit.Validate(); err != nil {
		return fmt.Errorf("App.KafkaRateLimit: %w", err)
	}
	if err := c.BinlogReader.Validate(); err != nil {
		return fmt.Errorf("App.BinlogReader: %w", err)
	}
	if c.Outbox == nil {
		c.Outbox = &Outbox{}
	}
	if err := c.Outbox.Validate(); err != nil {
		return fmt.Errorf("App.Outbox: %w", err)
	}
	if err := c.Breaker.Validate(); err != nil {
		return fmt.Errorf("App.Breaker: %w", err)
	}
	// OTEL is optional; validate only when explicitly enabled
	if c.OTEL != nil && c.OTEL.Enabled {
		c.OTEL.ApplyDefaults()
		if c.OTEL.Endpoint == "" {
			return fmt.Errorf("App.OTEL.Endpoint: %w", myErrors.ErrEmptyInput)
		}
		if c.OTEL.ServiceName == "" {
			return fmt.Errorf("App.OTEL.ServiceName: %w", myErrors.ErrEmptyInput)
		}
		if c.OTEL.TraceSampleRatio < 0 || c.OTEL.TraceSampleRatio > 1 {
			return fmt.Errorf("App.OTEL.TraceSampleRatio: %w: must be in [0,1], got %.4f", myErrors.ErrInvalidNumber, c.OTEL.TraceSampleRatio)
		}
		if c.OTEL.MetricExportInterval <= 0 {
			return fmt.Errorf("App.OTEL.MetricExportInterval: %w: must be positive, got %s", myErrors.ErrInvalidNumber, c.OTEL.MetricExportInterval)
		}
		if c.OTEL.MetricExportTimeout <= 0 {
			return fmt.Errorf("App.OTEL.MetricExportTimeout: %w: must be positive, got %s", myErrors.ErrInvalidNumber, c.OTEL.MetricExportTimeout)
		}
		if c.OTEL.MetricExportTimeout >= c.OTEL.MetricExportInterval {
			return fmt.Errorf("App.OTEL.MetricExportTimeout: %w: must be smaller than MetricExportInterval (%s >= %s)", myErrors.ErrInvalidNumber, c.OTEL.MetricExportTimeout, c.OTEL.MetricExportInterval)
		}
	}
	return nil
}
