package gorm

import (
	"context"
	"gochat/internal/shared/kernel"

	"gorm.io/gorm"
)

var _ kernel.UnitOfWork = (*UnitOfWork)(nil)

const unitOfWorkKey = "unit_of_work"

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Execute(ctx context.Context, fn func(context.Context) error) error {
	return u.DB(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, unitOfWorkKey, tx)
		return fn(ctxWithTx)
	})
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
