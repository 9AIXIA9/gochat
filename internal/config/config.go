package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gochat/internal/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultLanguage = "zh"
)

// Config 应用程序配置结构
type Config struct {
	Name      string        `mapstructure:"Name"`
	Host      string        `mapstructure:"Host"`
	Port      int           `mapstructure:"Port"`
	Language  string        `mapstructure:"Language"`
	Timeout   time.Duration `mapstructure:"Timeout"`
	Cert      *Cert         `mapstructure:"Cert"`
	Cookie    *Cookie       `mapstructure:"Cookie"`
	CORS      *CORS         `mapstructure:"CORS"`
	Token     *Token        `mapstructure:"Token"`
	RateLimit *RateLimit    `mapstructure:"RateLimit"`
	Database  *Database     `mapstructure:"Database"`
	Redis     *Redis        `mapstructure:"Redis"`
	Log       *Log          `mapstructure:"Log"`
	Snowflake *Snowflake    `mapstructure:"Snowflake"`
}

type Cert struct {
	HTTPSKeyFile  string `mapstructure:"HTTPSKeyFile"`
	HTTPSCertFile string `mapstructure:"HTTPSCertFile"`
}

type CORS struct {
	Origins          []string      `mapstructure:"Origins"`
	AllowMethods     []string      `mapstructure:"AllowMethods"`
	AllowHeaders     []string      `mapstructure:"AllowHeaders"`
	ExposeHeaders    []string      `mapstructure:"ExposeHeaders"`
	AllowCredentials bool          `mapstructure:"AllowCredentials"`
	MaxAge           time.Duration `mapstructure:"MaxAge"`
}

type Cookie struct {
	Domain   string `mapstructure:"Domain"`
	Secure   bool   `mapstructure:"Secure"`
	HttpOnly bool   `mapstructure:"HttpOnly"`
	Path     string `mapstructure:"Path"`
}

type Token struct {
	Refresh *RefreshToken
	Auth    *AuthToken
}

type RefreshToken struct {
	Length         int           `mapstructure:"Length"`
	ExpireDuration time.Duration `mapstructure:"ExpireDuration"`
}

type AuthToken struct {
	Secret         string        `mapstructure:"Secret"`
	ExpireDuration time.Duration `mapstructure:"ExpireDuration"`
}

type RateLimit struct {
	Period time.Duration `mapstructure:"Period"`
	Limit  int64         `mapstructure:"Limit"`
}

type Log struct {
	Mode       string `mapstructure:"Mode"`
	Level      string `mapstructure:"Level"`
	Filename   string `mapstructure:"Filename"`
	MaxSize    int    `mapstructure:"MaxSize"`
	MaxAge     int    `mapstructure:"MaxAge"`
	MaxBackups int    `mapstructure:"MaxBackups"`
}

type Snowflake struct {
	Node int64 `mapstructure:"Node"`
}

type Database struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
	Database string `mapstructure:"Database"`
}

type Redis struct {
	Host     string `mapstructure:"Host"`
	Port     int    `mapstructure:"Port"`
	Password string `mapstructure:"Password"`
	DB       int    `mapstructure:"DB"`
}

// MustLoad 加载配置，如果失败则退出
func MustLoad(configFile string) *Config {
	config, err := Load(configFile)
	if err != nil {
		log.Fatalf("config load failed,err:%v", err)
	}
	return config
}

// Load 加载配置文件
func Load(configFile string) (*Config, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configFile)

	// 读取配置文件内容
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	// 扩展环境变量
	expandedContent := os.ExpandEnv(string(content))

	// 使用扩展后的内容配置Viper
	ext := filepath.Ext(configFile)
	if ext == "" {
		ext = ".yaml"
	}
	v.SetConfigType(strings.TrimPrefix(ext, "."))

	// 使用扩展后的内容配置Viper
	if err := v.ReadConfig(strings.NewReader(expandedContent)); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parse config file failed: %w", err)
	}

	//验证配置正确性
	cfg.Validate()
	return cfg, nil
}

func (c *Config) Validate() {
	if len(c.Language) == 0 {
		c.Language = defaultLanguage
	} else if c.Language != "zh" && c.Language != "en" {
		log.Fatalf("language can only be selected from Chinese (zh) and English (en).")
	}

	if c.Token.Refresh.Length <= 0 {
		log.Fatalf("refresh token can't <= 0")
	}

	if err := utils.ValidateAllSubStructsNotEmpty(c); err != nil {
		log.Fatalf("an empty pointer appears:err:%v", err)
	}
}
