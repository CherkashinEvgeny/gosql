package sql

import (
	"context"

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
}

func (n nestedPropagation) Apply(factory base.Factory) (newFactory base.Factory) {
	return nestedPropagationFactory{factory}
}

type nestedPropagationFactory struct {
	parent base.Factory
}

func (f nestedPropagationFactory) Executor() (executor any) {
	return f.parent.Executor()
}

func (f nestedPropagationFactory) Tx(ctx context.Context, tx base.Tx, options ...any) (newTx base.Tx, err error) {
	if tx == nil {
		return f.parent.Tx(ctx, tx, options...)
	}
	executor := tx.Executor().(Executor)
	// TODO: generate uuid
	id := ""
	_, err = executor.Exec("SAVEPOINT $1", id)
	if err != nil {
		return nil, err
	}
	return nestedPropagationTx{executor, id}, nil
}

type nestedPropagationTx struct {
	executor Executor
	id       string
}

func (n nestedPropagationTx) Executor() (executor any) {
	return n.executor
}

func (n nestedPropagationTx) Commit(_ context.Context) (err error) {
	_, err = n.executor.Exec("RELEASE SAVEPOINT  $1", n.id)
	return err
}

func (n nestedPropagationTx) Rollback(_ context.Context) {
	_, _ = n.executor.Exec("ROLLBACK TO SAVEPOINT $1", n.id)
}
