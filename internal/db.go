package internal

import "context"

type Db interface {
	Executor
	Tx(ctx context.Context, tx Tx, options Options) (newTx Tx, err error)
}
