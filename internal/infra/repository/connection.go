package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gochat/internal/config"
	"gochat/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	KeyPrefix = "gochat"
)

func ConnectToMysql(conf *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to mysql failed,err:%w", err)
	}

	if err := model.AutoMigrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToRedis(conf *config.Redis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connect to redis failed,err:%w", err)
	}

	return rdb, nil
}
