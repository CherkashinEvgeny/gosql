package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CherkashinEvgeny/gosql/internal"
)

type Manager struct {
	base internal.Manager
}

type Option = internal.Option

func New(db *sql.DB, options ...Option) Manager {
	return Manager{internal.New(
		fmt.Sprintf("%p", db),
		&baseDb{db, options},
	)}
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) error, options ...Option) (err error) {
	return m.base.Transactional(ctx, f, options...)
}

type BeginError = internal.BeginError

type CommitError = internal.CommitError

type RollbackError = internal.RollbackError

func (m Manager) Executor(ctx context.Context) (executor Executor) {
	return m.base.Executor(ctx).(Executor)
}
