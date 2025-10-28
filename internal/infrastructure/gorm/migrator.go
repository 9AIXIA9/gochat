package gorm

import "gorm.io/gorm"

type Model interface {
	TableName() string
}

func AutoMigrate(db *gorm.DB, models ...Model) error {
	var tables []interface{}
	for _, m := range models {
		tables = append(tables, m)
	}
	return db.AutoMigrate(tables...)
}
