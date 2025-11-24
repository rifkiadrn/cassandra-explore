package usecase

import "context"

type Transaction interface {
	Commit() error
	Rollback() error
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, context.Context, error)
}
