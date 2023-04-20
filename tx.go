package sql

import (
	"context"
	"database/sql"
)

type Tx struct {
	*sql.Tx
	*TxOptions
}

type TxOptions = sql.TxOptions

const (
	LevelDefault         = sql.LevelDefault
	LevelReadUncommitted = sql.LevelReadUncommitted
	LevelReadCommitted   = sql.LevelReadCommitted
	LevelWriteCommitted  = sql.LevelWriteCommitted
	LevelRepeatableRead  = sql.LevelRepeatableRead
	LevelSnapshot        = sql.LevelSnapshot
	LevelSerializable    = sql.LevelSerializable
	LevelLinearizable    = sql.LevelLinearizable
)

var _ Executor = (*Tx)(nil)

func (t *Tx) NamedExec(query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.Exec(query, posArgs...)
}

func (t *Tx) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.ExecContext(ctx, query, posArgs...)
}

func (t *Tx) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.Query(query, posArgs...)
}

func (t *Tx) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.QueryContext(ctx, query, posArgs...)
}

func (t *Tx) NamedQueryRow(query string, args map[string]any) (row *Row) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.QueryRow(query, posArgs...)
}

func (t *Tx) NamedQueryRowContext(ctx context.Context, query string, args map[string]any) (row *Row) {
	query, posArgs := convertToPositionalParameters(query, args)
	return t.QueryRowContext(ctx, query, posArgs...)
}
