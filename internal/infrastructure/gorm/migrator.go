package gorm

import (
	"fmt"

	"gorm.io/gorm"
)

type Table interface {
	TableName() string
}

func AutoMigrate(db *gorm.DB, models ...Table) error {
	if len(models) == 0 {
		return nil
	}
	tables := make([]interface{}, 0, len(models))

	for i, model := range models {
		if model == nil {
			return fmt.Errorf("model at index %d is nil", i)
		}

		if n := model.TableName(); n == "" {
			return fmt.Errorf("model at index %d (%T) has empty table name", i, model)
		}

		tables = append(tables, model)
	}

	return db.AutoMigrate(tables...)
}
