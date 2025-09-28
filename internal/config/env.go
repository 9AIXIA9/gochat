package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

//TODO 自动替换环境变量

// loadEnvFile 根据当前环境加载对应的环境变量文件
func loadEnvFile() error {
	env := os.Getenv("GOCHAT_ENV")
	if env == "" {
		env = "development" // 默认��境
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
		"CORS.Origins":      "CORS_ORIGINS",
		"Cookie.Domain":     "COOKIE_DOMAIN",
		"Cookie.Secure":     "COOKIE_SECURE",
		"Cookie.HttpOnly":   "COOKIE_HTTPONLY",
		"Cookie.Path":       "COOKIE_PATH",
	}

	for configPath, envVar := range envMappings {
		err := v.BindEnv(configPath, envVar)
		if err != nil {
			log.Fatalf("bind env failed,err:%v", err)
		}
	}

	// 特殊处理 CORS_ORIGINS：支持 JSON 数组或逗号分隔、或单个 origin
	if env := os.Getenv("CORS_ORIGINS"); env != "" {
		var origins []string
		// 尝试解析为 JSON 数组
		if err := json.Unmarshal([]byte(env), &origins); err == nil {
			v.Set("CORS.Origins", origins)
			return
		}
		// 回退为逗号分隔或单值，例如: https://a,https://b 或 https://a
		s := strings.TrimSpace(env)
		s = strings.Trim(s, "[]")
		parts := strings.Split(s, ",")
		for i, p := range parts {
			parts[i] = strings.Trim(strings.TrimSpace(p), "\"'")
		}
		// 过滤空字符串
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if p != "" {
				out = append(out, p)
			}
		}
		v.Set("CORS.Origins", out)
	}
}
