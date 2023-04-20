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
	BeginTx(ctx context.Context, options *TxOptions) (txCtx context.Context, err error)
	EndTx(ctx context.Context, txErr error) (rollbackErr error, commitErr error)
	CommitTx(ctx context.Context) (err error)
	RollbackTx(ctx context.Context) (err error)
}
