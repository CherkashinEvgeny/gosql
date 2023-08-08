package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CherkashinEvgeny/gosql/base"
)

func Isolation(level sql.IsolationLevel) (option Option) {
	return isolationOption{level}
}

type isolationOption struct {
	level sql.IsolationLevel
}

func (o isolationOption) Name() (name string) {
	return "isolation"
}

func (o isolationOption) Priority() (priority int) {
	return 5
}

func (o isolationOption) Apply(factory base.Factory) (newFactory base.Factory) {
	return isolationOptionFactory{factory, o.level}
}

type isolationOptionFactory struct {
	next  base.Factory
	level sql.IsolationLevel
}

const isolationKey = "isolation"

func (f isolationOptionFactory) Tx(ctx context.Context, tx base.Tx, values base.Valuer) (newTx base.Tx, err error) {
	if tx == nil {
		return f.next.Tx(ctx, tx, base.WithValue(values, isolationKey, f.level))
	}
	level, ok := tx.Value(isolationKey).(sql.IsolationLevel)
	if ok && f.level > level {
		return nil, fmt.Errorf("isolation level %s is to low", f.level)
	}
	return f.next.Tx(ctx, tx, base.WithValue(values, isolationKey, f.level))
}
