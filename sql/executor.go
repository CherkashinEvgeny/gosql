package sql

import (
	"context"
	"database/sql"
)

type Executor interface {
	Exec(query string, args ...any) (result sql.Result, err error)
	ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error)
	Query(query string, args ...any) (rows *sql.Rows, err error)
	QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error)
	QueryRow(query string, args ...any) (row *sql.Row)
	QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row)
}
