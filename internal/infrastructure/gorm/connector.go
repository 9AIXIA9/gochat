package gorm

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
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
		defer func() {
			if db == nil {
				return
			}
			sqlDB, err := db.DB()
			if err != nil {
				zap.L().Warn("close mysql connection failed", zap.Error(err))
				return
			}
			if err := sqlDB.Close(); err != nil {
				zap.L().Warn("close mysql connection failed", zap.Error(err))
			}
		}()
		return nil, fmt.Errorf("connect to mysql failed,err:%w", err)
	}
	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, fmt.Errorf("register gorm otel plugin failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get mysql sql.DB failed: %w", err)
	}
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetimeSec > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetimeSec) * time.Second)
	}
	if config.ConnMaxIdleTimeSec > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTimeSec) * time.Second)
	}

	return db, nil
}
