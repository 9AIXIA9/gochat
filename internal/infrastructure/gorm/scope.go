package gorm

import "gorm.io/gorm"

// Scope 抽象，便于组合查询
type Scope func(*gorm.DB) *gorm.DB

func applyScopes(db *gorm.DB, scopes ...Scope) *gorm.DB {
	for _, s := range scopes {
		db = s(db)
	}
	return db
}

func Where(query interface{}, args ...interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB { return db.Where(query, args...) }
}
func Preload(association string, conds ...interface{}) Scope {
	return func(db *gorm.DB) *gorm.DB { return db.Preload(association, conds...) }
}
func Order(order string) Scope {
	return func(db *gorm.DB) *gorm.DB { return db.Order(order) }
}

func Paginate(offset, limit int) Scope {
	return func(db *gorm.DB) *gorm.DB {
		if offset < 0 {
			offset = 0
		}
		db = db.Offset(offset)
		if limit > 0 {
			db = db.Limit(limit)
		}
		return db
	}
}
