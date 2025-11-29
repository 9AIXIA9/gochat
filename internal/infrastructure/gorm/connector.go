package gorm

import (
	"fmt"

	"go.opentelemetry.io/otel"
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
		return nil, fmt.Errorf("connect to mysql failed,err:%w", err)
	}
	// Register tracing plugin only if a tracer provider exists (telemetry may be disabled)
	if otel.GetTracerProvider() != nil {
		_ = db.Use(tracing.NewPlugin())
	}
	// Register custom observability plugin (metrics + slow query tracing)
	metrics := NewDBMetrics(nil)
	metrics.Register(nil)
	_ = db.Use(NewObservabilityPlugin(metrics, config.SlowThresholdMillis))
	return db, nil
}
