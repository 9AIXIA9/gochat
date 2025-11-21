package gorm

import (
	"context"
	"time"

	"gochat/internal/infrastructure/prometheus"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ kernel.UnitOfWork = (*UnitOfWork)(nil)

const unitOfWorkKey = "unit_of_work"

//TODO 实现不妥

type UnitOfWork struct {
	db      *gorm.DB
	metrics *prometheus.Metrics
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

// SetMetrics attaches metrics collectors (optional).
func (u *UnitOfWork) SetMetrics(m *prometheus.Metrics) { u.metrics = m }

func (u *UnitOfWork) Execute(ctx context.Context, fn func(context.Context) error) error {
	start := time.Now()
	err := u.DB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, unitOfWorkKey, tx)
		return fn(ctxWithTx)
	})
	if u.metrics != nil {
		res := "success"
		if err != nil {
			res = "error"
		}
		u.metrics.UOWDur.WithLabelValues(res).Observe(time.Since(start).Seconds())
		u.metrics.UOWTotal.WithLabelValues(res).Inc()
	}
	return err
}

func (u *UnitOfWork) DB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return u.db
	}
	txInCtx := ctx.Value(unitOfWorkKey)
	tx, ok := txInCtx.(*gorm.DB)
	if !ok {
		return u.db
	}
	return tx
}
