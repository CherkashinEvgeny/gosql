package base

import "context"

type Option interface {
	Apply(factory Factory) (newFactory Factory)
}

type Factory interface {
	Executor() (executor any)
	Tx(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error)
}

type Tx interface {
	Executor() (executor any)
	Commit(ctx context.Context) (err error)
	Rollback(ctx context.Context)
}

type Store interface {
	Get(ctx context.Context) (core Tx)
	Set(ctx context.Context, core Tx) (newCtx context.Context)
}
