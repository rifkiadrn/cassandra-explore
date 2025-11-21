package context_db

import (
	"context"

	"gorm.io/gorm"
)

type key struct{}

var dbKey key

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, tx)
}

// BeginTxWithContext starts a transaction and returns the new context and tx object
func BeginTxWithContext(ctx context.Context, db *gorm.DB) (*gorm.DB, context.Context) {
	tx := db.WithContext(ctx).Begin()
	txCtx := context.WithValue(ctx, dbKey, tx)
	return tx, txCtx
}

func GetTx(ctx context.Context) (*gorm.DB, error) {
	db, ok := ctx.Value(dbKey).(*gorm.DB)
	if !ok || db == nil {
		return nil, gorm.ErrInvalidDB
	}
	return db, nil
}
