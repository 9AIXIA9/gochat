package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

// LoadEnvFile 根据当前环境加载对应的环境变量文件
func LoadEnvFile() error {
	env := os.Getenv("GOCHAT_ENV")
	if env == "" {
		env = "development" // 默认环境
	}

	// 首先尝试加载特定环境的配置
	envFile := ".env." + env
	err := godotenv.Load(envFile)

	// 如果特定环境配置不存在，尝试加载默认配置
	if err != nil && !os.IsNotExist(err) {
		// 不存在文件是正常的，其他错误才需要报告
		return fmt.Errorf("加载环境文件 %s 失败: %w", envFile, err)
	}

	return nil
}
