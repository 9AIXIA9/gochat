package repository

import (
	"fmt"
	"gochat/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func MustConnectToMysql(database config.Database) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		database.Username,
		database.Password,
		database.Host,
		database.Port,
		database.Database)

	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect to database failed,err:%v", err)
	}

	return db
}
