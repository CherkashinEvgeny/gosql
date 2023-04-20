package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
)

type DB struct {
	*sql.DB
	txKey string
}

func Open(driverName string, dataSource string) (db *DB, err error) {
	sqlDb, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	db = &DB{DB: sqlDb}
	db.txKey = txKey(db)
	return db, nil
}

func OpenDB(conn driver.Connector) (db *DB) {
	sqlDb := sql.OpenDB(conn)
	db = &DB{DB: sqlDb}
	db.txKey = txKey(db)
	return db
}

func (d *DB) Exec(query string, arg ...any) (result Result, err error) {
	return d.DB.Exec(query, arg...)
}

func (d *DB) NamedExec(query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Exec(query, posArgs...)
}

func (d *DB) ExecContext(ctx context.Context, query string, arg ...any) (result Result, err error) {
	return d.DB.ExecContext(ctx, query, arg...)
}

func (d *DB) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.ExecContext(ctx, query, posArgs...)
}

func (d *DB) Query(query string, args ...any) (rows *Rows, err error) {
	return d.DB.Query(query, args...)
}

func (d *DB) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Query(query, posArgs...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (rows *Rows, err error) {
	return d.DB.QueryContext(ctx, query, args...)
}

func (d *DB) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.QueryContext(ctx, query, posArgs...)
}

var _ Executor = (*DB)(nil)

func (d *DB) Executor(ctx context.Context) (executor Executor) {
	tx, found := d.extractTxFromContext(ctx)
	if found {
		executor = tx
	} else {
		executor = d
	}
	return executor
}

var _ TxController = (*DB)(nil)

func (d *DB) Transactional(
	ctx context.Context,
	f func(ctx context.Context) error,
	options *TxOptions,
) (beginErr error, err error, commitErr error, rollbackErr error) {
	ctx, beginErr = d.Begin(ctx, options)
	if beginErr != nil {
		return beginErr, nil, nil, nil
	}
	defer func() {
		commitErr, rollbackErr = d.End(ctx, err)
	}()
	err = f(ctx)
	return beginErr, err, nil, nil
}

func (d *DB) Begin(ctx context.Context, options *TxOptions) (txCtx context.Context, err error) {
	tx, found := d.extractTxFromContext(ctx)
	if found {
		if tx.Isolation != sql.LevelDefault && tx.Isolation < options.Isolation {
			err = transactionIsolationLevelIsToLow
			cfg.LogicalErrorHandler(err)
			return ctx, err
		}
		return ctx, nil
	}
	sqlTx, err := d.DB.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}
	tx = &Tx{sqlTx, options}
	return context.WithValue(ctx, d.txKey, tx), nil
}

var transactionIsolationLevelIsToLow = errors.New("transaction isolation level is to low")

func (d *DB) End(ctx context.Context, txErr error) (commitErr error, rollbackErr error) {
	if txErr == nil {
		commitErr = d.Commit(ctx)
	} else {
		rollbackErr = d.Rollback(ctx)
	}
	return commitErr, rollbackErr
}

func (d *DB) Commit(ctx context.Context) (err error) {
	tx, found := d.extractTxFromContext(ctx)
	if !found {
		err = transactionNotFound
		cfg.LogicalErrorHandler(err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		if err == sql.ErrTxDone {
			cfg.LogicalErrorHandler(err)
		}
		return err
	}
	return err
}

func (d *DB) Rollback(ctx context.Context) (err error) {
	tx, found := d.extractTxFromContext(ctx)
	if !found {
		err = transactionNotFound
		cfg.LogicalErrorHandler(err)
		return err
	}
	err = tx.Rollback()
	if err != nil {
		if err == sql.ErrTxDone {
			cfg.LogicalErrorHandler(err)
		}
		return err
	}
	return nil
}

var transactionNotFound = errors.New("transaction not found")

func (d *DB) extractTxFromContext(ctx context.Context) (tx *Tx, found bool) {
	txValue := ctx.Value(d.txKey)
	if txValue == nil {
		return nil, false
	}
	var txCastOk bool
	tx, txCastOk = txValue.(*Tx)
	if !txCastOk {
		return nil, false
	}
	return tx, true
}

func txKey(d *DB) (key string) {
	return fmt.Sprintf("tx-%p", d)
}
