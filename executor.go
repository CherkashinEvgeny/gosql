package sqlx

import (
	"context"
	"database/sql"
)

type Executor interface {
	Exec(query string, args ...any) (result Result, err error)
	NamedExec(query string, args map[string]any) (result Result, err error)

	ExecContext(ctx context.Context, query string, args ...any) (result Result, err error)
	NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error)

	Query(query string, args ...any) (rows *Rows, err error)
	NamedQuery(query string, args map[string]any) (rows *Rows, err error)

	QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error)
	NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error)
}

type Result = sql.Result

type Rows = sql.Rows

var _ Executor = (*DB)(nil)

func (d *DB) Executor(ctx context.Context) (executor Executor) {
	tx, err := d.extractTxFromContext(ctx)
	if err == nil {
		executor = tx
	} else {
		executor = d
	}
	return executor
}

func (d *DB) Exec(query string, arg ...any) (result Result, err error) {
	return (*sql.DB)(d).Exec(query, arg...)
}

func (d *DB) NamedExec(query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Exec(query, posArgs...)
}

func (d *DB) ExecContext(ctx context.Context, query string, arg ...any) (result Result, err error) {
	return (*sql.DB)(d).ExecContext(ctx, query, arg...)
}

func (d *DB) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.ExecContext(ctx, query, posArgs...)
}

func (d *DB) Query(query string, args ...any) (rows *Rows, err error) {
	return (*sql.DB)(d).Query(query, args...)
}

func (d *DB) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Query(query, posArgs...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error) {
	return (*sql.DB)(d).QueryContext(ctx, query, args...)
}

func (d *DB) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.QueryContext(ctx, query, posArgs...)
}

var _ Executor = (*Tx)(nil)

func (t *Tx) Exec(query string, arg ...any) (result Result, err error) {
	return (*sql.Tx)(t).Exec(query, arg...)
}

func (t *Tx) NamedExec(query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.Exec(query, posArgs...)
}

func (t *Tx) ExecContext(ctx context.Context, query string, arg ...any) (result Result, err error) {
	return (*sql.Tx)(t).ExecContext(ctx, query, arg...)
}

func (t *Tx) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.ExecContext(ctx, query, posArgs...)
}

func (t *Tx) Query(query string, args ...any) (rows *Rows, err error) {
	return (*sql.Tx)(t).Query(query, args...)
}

func (t *Tx) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.Query(query, posArgs...)
}

func (t *Tx) QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error) {
	return (*sql.Tx)(t).QueryContext(ctx, query, args...)
}

func (t *Tx) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.QueryContext(ctx, query, posArgs...)
}

func convertToPositionalParameters(query string, args map[string]any) (string, []any) {
	// TODO
	return query, nil
}
