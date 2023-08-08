package sql

import (
	"context"
	"sync/atomic"

	"github.com/CherkashinEvgeny/gosql/base"
)

var (
	Required  Option = base.Required
	Supports  Option = base.Supports
	Mandatory Option = base.Mandatory
	Never     Option = base.Never
	Nested    Option = nestedPropagation{}
)

type nestedPropagation struct {
	base.Propagation
	id int64
}

func (n nestedPropagation) Apply(factory base.Factory) (newFactory base.Factory) {
	return nestedPropagationFactory{factory, &n.id}
}

type nestedPropagationFactory struct {
	next base.Factory
	id   *int64
}

func (f nestedPropagationFactory) Tx(ctx context.Context, tx base.Tx, options base.Valuer) (newTx base.Tx, err error) {
	if tx == nil {
		return f.next.Tx(ctx, tx, options)
	}
	executor := tx.Executor().(Executor)
	id := atomic.AddInt64(f.id, 1)
	_, err = executor.Exec("SAVEPOINT $1", id)
	if err != nil {
		return nil, err
	}
	return nestedPropagationTx{tx, id}, nil
}

type nestedPropagationTx struct {
	base.Tx
	id int64
}

func (n nestedPropagationTx) Commit(ctx context.Context) (err error) {
	_, err = n.Executor().(baseTx).tx.ExecContext(ctx, "RELEASE SAVEPOINT $1", n.id)
	return err
}

func (n nestedPropagationTx) Rollback(ctx context.Context) {
	_, _ = n.Executor().(baseTx).tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT $1", n.id)
}
