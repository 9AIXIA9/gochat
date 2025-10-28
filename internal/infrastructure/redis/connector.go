package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(config *Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connect to redis failed,err:%w", err)
	}

	return rdb, nil
}
