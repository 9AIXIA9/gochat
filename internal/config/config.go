package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name      string    `yaml:"Name"`
	Host      string    `yaml:"Host"`
	Port      int       `yaml:"Port"`
	JWT       JWT       `yaml:"JWT"`
	Database  Database  `yaml:"Database"`
	Redis     Redis     `yaml:"Redis"`
	Log       Log       `yaml:"Log"`
	Snowflake Snowflake `yaml:"Snowflake"`
}

type JWT struct {
	Secret     string `yaml:"Secret"`
	ExpireTime int    `yaml:"ExpireTime"` // hours
}

type Log struct {
	Level      string `yaml:"Level"`
	Filename   string `yaml:"Filename"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxAge     int    `yaml:"MaxAge"`
	MaxBackups int    `yaml:"MaxBackups"`
}

type Snowflake struct {
	Node int64 `yaml:"Node"`
}

type Database struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

type Redis struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

func Init(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	c := &Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	return c, nil
}
