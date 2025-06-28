package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

const (
	defaultLanguage = "zh"
)

// Config 应用程序配置结构
type Config struct {
	Name      string     `mapstructure:"Name"`
	Host      string     `mapstructure:"Host"`
	Port      int        `mapstructure:"Port"`
	Language  string     `mapstructure:"Language"`
	JWT       *JWT       `mapstructure:"JWT"`
	Database  *Database  `mapstructure:"Database"`
	Redis     *Redis     `mapstructure:"Redis"`
	Log       *Log       `mapstructure:"Log"`
	Snowflake *Snowflake `mapstructure:"Snowflake"`
}

type JWT struct {
	Secret     string `mapstructure:"Secret"`
	ExpireTime int    `mapstructure:"ExpireTime"` // 单位：小时
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

// MustLoad 加载配置，如果失败则panic
func MustLoad(configFile string) *Config {
	config, err := Load(configFile)
	if err != nil {
		panic(err)
	}
	return config
}

// Load 加载配置文件和环境变量
func Load(configFile string) (*Config, error) {
	// 加载环境变量
	if err := loadEnvFile(); err != nil {
		return nil, fmt.Errorf("加载环境变量文件失败: %w", err)
	}

	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(configFile)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 配置环境变量设置
	setupEnvironmentVars(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() {
	if len(c.Language) == 0 {
		c.Language = defaultLanguage
	} else if c.Language != "zh" && c.Language != "en" {
		log.Fatalf("language can only be selected from Chinese (zh) and English (en).")
	}
}
