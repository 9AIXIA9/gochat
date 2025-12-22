package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9" // 添加这行
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectToRedis(config *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.User,
		Password: config.Password,
		DB:       config.Database,
	})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, fmt.Errorf("instrument redis client failed, err:%w", err)
	}

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		defer func() {
			if err := rdb.Close(); err != nil {
				zap.L().Warn("close redis client failed", zap.Error(err))
			}
		}()
		return nil, fmt.Errorf("connect to redis failed,err:%w", err)
	}

	return rdb, nil
}
