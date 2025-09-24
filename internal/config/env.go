package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// loadEnvFile 根据当前环境加载对应的环境变量文件
func loadEnvFile() error {
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

// setupEnvironmentVars 配置环境变量处理
func setupEnvironmentVars(v *viper.Viper) {
	// 设置环境变量前缀和替换规则
	v.SetEnvPrefix("GOCHAT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 绑定特定环境变量到配置路径
	envMappings := map[string]string{
		"Database.Password": "DB_PASSWORD",
		"Token.Secret":      "JWT_SECRET",
		"Redis.Password":    "REDIS_PASSWORD",
	}

	for configPath, envVar := range envMappings {
		err := v.BindEnv(configPath, envVar)
		if err != nil {
			log.Fatalf("bind env failed,err:%v", err)
		}
	}
}
