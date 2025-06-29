package repository

import (
	"fmt"
	"gochat/internal/config"
	"gochat/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func MustConnectToMysql(database config.Database) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		database.Username,
		database.Password,
		database.Host,
		database.Port,
		database.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect to database failed,err:%v", err)
	}

	model.AutoMigrate(db)

	return db
}
