package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

// loadEnvFile 根据当前环境加载对应的环境变量文件
func loadEnvFile(env string) error {
	// 首先尝试加载特定环境的配置
	envFile := ".env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("load env file %s failed: %w", envFile, err)
	}
	log.Printf("using environment: %s", env)
	return nil
}
