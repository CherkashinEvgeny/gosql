package sql

import (
	"context"
)

type TxController interface {
	Transactional(
		ctx context.Context,
		f func(ctx context.Context) error,
		options *TxOptions,
	) (beginErr error, err error, commitErr error, rollbackErr error)
	Begin(ctx context.Context, options *TxOptions) (txCtx context.Context, err error)
	End(ctx context.Context, txErr error) (rollbackErr error, commitErr error)
	Commit(ctx context.Context) (err error)
	Rollback(ctx context.Context) (err error)
}
