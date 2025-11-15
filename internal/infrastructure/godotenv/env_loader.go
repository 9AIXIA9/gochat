package godotenv

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvFile(path string) error {
	// 优先加载指定的环境文件
	if _, err := os.Stat(path); err == nil {
		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("load env file %s failed: %w", path, err)
		}
		log.Printf("using environment: %s", path)
		return nil
	}

	// 回退尝试加载通用 .env
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return fmt.Errorf("load env file .env failed: %w", err)
		}
		log.Printf("using environment: .env)")
		return nil
	}

	// 两者都不存在则不报错，允许使用进程环境变量（例如由 docker-compose 注入）
	log.Printf("using environment: process env (no .env files found)")
	return nil
}
