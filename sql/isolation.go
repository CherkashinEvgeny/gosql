package sql

import (
	"context"
	"database/sql"

	"github.com/CherkashinEvgeny/gosql/base"
)

func Isolation(level sql.IsolationLevel) (option base.Option) {
	return isolationOption(level)
}

type isolationOption sql.IsolationLevel

func (i isolationOption) Apply(factory base.Factory) (newFactory base.Factory) {
	return isolationFactory{factory, sql.IsolationLevel(i)}
}

type isolationFactory struct {
	parent    base.Factory
	isolation sql.IsolationLevel
}

func (f isolationFactory) Executor() (executor any) {
	return f.parent.Executor()
}

func (f isolationFactory) Tx(ctx context.Context, tx base.Tx, options ...any) (newTx base.Tx, err error) {
	options = append(options, f.isolation)
	return f.parent.Tx(ctx, tx, options...)
}
