package context_db

import (
	"context"

	"github.com/rifkiadrn/cassandra-explore/internal/usecase"
	"gorm.io/gorm"
)

type GormTransaction struct {
	tx *gorm.DB
}

func (t *GormTransaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *GormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

type GormUnitOfWork struct {
	db *gorm.DB
}

func NewGormUnitOfWork(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{db: db}
}

type txKey struct{}

func (u *GormUnitOfWork) Begin(ctx context.Context) (usecase.Transaction, context.Context, error) {
	tx := u.db.Begin()

	// store tx in context
	txCtx := context.WithValue(ctx, txKey{}, tx)

	return &GormTransaction{tx}, txCtx, nil
}

func GetTx(ctx context.Context) *gorm.DB {
	tx := ctx.Value(txKey{})
	if tx == nil {
		return nil
	}
	return tx.(*gorm.DB)
}
