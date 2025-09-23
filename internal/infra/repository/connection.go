package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gochat/internal/config"
	"gochat/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

const (
	KeyPrefix = "gochat"
)

func MustConnectToMysql(conf *config.Database) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect to mysql failed,err:%v", err)
	}

	model.AutoMigrate(db)

	return db
}

func MustConnectToRedis(conf *config.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("connect to redis failed,err:%v", err)
	}

	return rdb
}
