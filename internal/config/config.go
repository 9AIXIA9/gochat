package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name      string     `yaml:"Name"`
	Host      string     `yaml:"Host"`
	Port      int        `yaml:"Port"`
	JWT       *JWT       `yaml:"JWT"`
	Database  *Database  `yaml:"Database"`
	Redis     *Redis     `yaml:"Redis"`
	Log       *Log       `yaml:"Log"`
	Snowflake *Snowflake `yaml:"Snowflake"`
}

type JWT struct {
	Secret     string `yaml:"Secret"`
	ExpireTime int    `yaml:"ExpireTime"` // hours
}

type Log struct {
	Mode       string `yaml:"Mode"`
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

func MustLoad(configFile string) *Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("read config file failed: %v", err)
		return nil
	}

	c := &Config{}
	if err := yaml.Unmarshal(data, c); err != nil {
		log.Fatalf("parse config failed: %v", err)
		return nil
	}

	return c
}
