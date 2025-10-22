package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToMysql(config *MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to mysql failed,err:%w", err)
	}
	return db, nil
}
