package gorm

import (
	"context"

	"gorm.io/gorm"
)

const transactionKey = "transaction"

func SetTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

func GetTransaction(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	txInCtx := ctx.Value(transactionKey)
	tx, ok := txInCtx.(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}
