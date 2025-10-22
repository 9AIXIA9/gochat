package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gochat/internal/shared/config"
	"os"
	"path/filepath"
	"strings"
)

func loadConfigFile(path string) (*config.App, error) {
	v := viper.New()

	// 设置配置文件
	v.SetConfigFile(path)

	// 读取配置文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	// 扩展环境变量
	expandedContent := os.ExpandEnv(string(content))

	// 使用扩展后的内容配置Viper
	ext := filepath.Ext(path)
	if ext == "" {
		ext = ".yaml"
	}
	v.SetConfigType(strings.TrimPrefix(ext, "."))

	// 使用扩展后的内容配置Viper
	if err := v.ReadConfig(strings.NewReader(expandedContent)); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	cfg := new(config.App)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parse config file failed: %w", err)
	}

	//验证配置正确性
	return cfg, cfg.Validate()
}
