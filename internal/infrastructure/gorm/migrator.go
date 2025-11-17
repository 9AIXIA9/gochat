package gorm

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, models ...interface{ TableName() string }) error {
	for _, model := range models {
		if n := model.TableName(); n == "" {
			return fmt.Errorf("model has empty table name")
		} else {
			if db.Migrator().HasTable(n) {
				zap.L().Info("table already exists, skipping migration", zap.String("table", n))
				return nil
			}
		}
	}

	var tables []interface{}
	for _, m := range models {
		tables = append(tables, m)
	}
	return db.AutoMigrate(tables...)
}
