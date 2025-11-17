package gorm

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB, models ...interface{ TableName() string }) error {
	var tables []interface{}
	for _, m := range models {
		tables = append(tables, m)
	}
	return db.AutoMigrate(tables...)
}

// AutoMigrateIfFresh checks whether the first model's table exists. If not, it treats the schema as fresh
// and runs AutoMigrate for all provided models. Returns a boolean indicating whether migration executed.
// This provides a simple "first run" behavior without needing an explicit flag.
func AutoMigrateIfFresh(db *gorm.DB, models ...interface{ TableName() string }) (bool, error) {
	if len(models) == 0 {
		return false, nil
	}
	// Use first model as sentinel. If absent, we assume a fresh database.
	sentinel := models[0].TableName()
	if db.Migrator().HasTable(sentinel) {
		return false, nil
	}
	if err := AutoMigrate(db, models...); err != nil {
		return false, err
	}
	return true, nil
}
