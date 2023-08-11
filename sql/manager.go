package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CherkashinEvgeny/gosql/internal"
)

type Executor interface {
	Exec(query string, args ...any) (result sql.Result, err error)
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	Query(query string, args ...any) (rows *sql.Rows, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRow(query string, args ...any) (row *sql.Row)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}

type Manager struct {
	base internal.Manager
}

type Option = internal.Option

func New(db *sql.DB, options ...Option) Manager {
	options = append(options, WithIsolationLevel(sql.LevelReadCommitted))
	return Manager{internal.New(
		fmt.Sprintf("%p", db),
		&baseDb{db},
	)}
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) error, options ...Option) (err error) {
	return m.base.Transactional(ctx, f, options...)
}

type BeginError = internal.BeginError

type CommitError = internal.CommitError

func (m Manager) Executor(ctx context.Context) (executor Executor) {
	return m.base.Executor(ctx).(Executor)
}
