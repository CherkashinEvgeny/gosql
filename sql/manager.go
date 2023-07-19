package sql

import (
	"context"
	"database/sql"

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

type Manager struct {
	base base.Manager
}

func New(db *sql.DB) Manager {
	return Manager{base.New(factory{db})}
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) error, options ...Option) (err error) {
	return m.base.Transactional(ctx, f, options...)
}

func (m Manager) Get(ctx context.Context) Executor {
	return m.base.Get(ctx).(Executor)
}

type factory struct {
	db *sql.DB
}

func (d factory) Executor() (executor any) {
	return d.db
}

func (d factory) Tx(ctx context.Context, _ base.Tx, options ...any) (newTx base.Tx, err error) {
	sqlOptions := &sql.TxOptions{}
	for _, option := range options {
		if level, ok := option.(sql.IsolationLevel); ok {
			sqlOptions.Isolation = level
		}
	}
	sqlTx, err := d.db.BeginTx(ctx, sqlOptions)
	if err != nil {
		return nil, err
	}
	return txWrapper{sqlTx}, nil
}

type txWrapper struct {
	tx *sql.Tx
}

func (t txWrapper) Executor() (executor any) {
	return t.tx
}

func (t txWrapper) Commit(ctx context.Context) (err error) {
	// TODO implement me
	panic("implement me")
}

func (t txWrapper) Rollback(ctx context.Context) {
	// TODO implement me
	panic("implement me")
}
