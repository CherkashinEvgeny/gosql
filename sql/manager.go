package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CherkashinEvgeny/gosql/base"
)

type Executor interface {
	Exec(query string, args ...any) (result sql.Result, err error)
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	Query(query string, args ...any) (rows *sql.Rows, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRow(query string, args ...any) (row *sql.Row)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}

type Option = base.Option

type BeginError = base.BeginError

type CommitError = base.CommitError

type Manager struct {
	base base.Manager
}

func New(db *sql.DB) Manager {
	return Manager{base.New(
		fmt.Sprintf("%p", db),
		(*baseDb)(db),
		(*baseFactory)(db),
	)}
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) error, options ...Option) (err error) {
	return m.base.Transactional(ctx, f, options...)
}

func (m Manager) Executor(ctx context.Context) (executor Executor) {
	return m.base.Executor(ctx).(Executor)
}

type baseDb sql.DB

func (d *baseDb) Executor() (executor any) {
	return (*sql.DB)(d)
}

type baseFactory sql.DB

func (d *baseFactory) Tx(ctx context.Context, _ base.Tx, options ...any) (newTx base.Tx, err error) {
	sqlOptions := &sql.TxOptions{}
	for _, option := range options {
		if level, ok := option.(sql.IsolationLevel); ok {
			sqlOptions.Isolation = level
		}
	}
	sqlTx, err := (*sql.DB)(d).BeginTx(ctx, sqlOptions)
	if err != nil {
		return nil, err
	}
	return baseTx{sqlTx}, nil
}

type baseTx struct {
	tx *sql.Tx
}

func (t baseTx) Executor() (executor any) {
	return t.tx
}

func (t baseTx) Commit(_ context.Context) (err error) {
	return t.tx.Commit()
}

func (t baseTx) Rollback(_ context.Context) {
	_ = t.tx.Rollback()
}
