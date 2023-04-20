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

func txKey(d *DB) (key string) {
	return fmt.Sprintf("tx-%p", d)
}

func (d *DB) Executor(ctx context.Context) (executor Executor) {
	tx, found := d.TxFromContext(ctx)
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
	ctx, beginErr = d.BeginTx(ctx, options)
	if beginErr != nil {
		return beginErr, nil, nil, nil
	}
	defer func() {
		commitErr, rollbackErr = d.EndTx(ctx, err)
	}()
	err = f(ctx)
	return beginErr, err, nil, nil
}

func (d *DB) BeginTx(ctx context.Context, options *TxOptions) (txCtx context.Context, err error) {
	tx, found := d.TxFromContext(ctx)
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

func (d *DB) EndTx(ctx context.Context, txErr error) (commitErr error, rollbackErr error) {
	if txErr == nil {
		commitErr = d.CommitTx(ctx)
	} else {
		rollbackErr = d.RollbackTx(ctx)
	}
	return commitErr, rollbackErr
}

func (d *DB) CommitTx(ctx context.Context) (err error) {
	tx, found := d.TxFromContext(ctx)
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

func (d *DB) RollbackTx(ctx context.Context) (err error) {
	tx, found := d.TxFromContext(ctx)
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

func (d *DB) TxFromContext(ctx context.Context) (tx *Tx, found bool) {
	txValue := ctx.Value(d.txKey)
	if txValue == nil {
		return nil, false
	}
	tx, txCastOk := txValue.(*Tx)
	if !txCastOk {
		return nil, false
	}
	return tx, true
}

var _ Executor = (*DB)(nil)

func (d *DB) NamedExec(query string, args map[string]any) (result Result, err error) {

	query, posArgs := convertToPositionalParameters(query, args)
	return d.Exec(query, posArgs...)
}

func (d *DB) NamedExecContext(ctx context.Context, query string, args map[string]any) (result Result, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.ExecContext(ctx, query, posArgs...)
}

func (d *DB) NamedQuery(query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.Query(query, posArgs...)
}

func (d *DB) NamedQueryContext(ctx context.Context, query string, args map[string]any) (rows *Rows, err error) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.QueryContext(ctx, query, posArgs...)
}

func (d *DB) NamedQueryRow(query string, args map[string]any) (row *Row) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.QueryRow(query, posArgs...)
}

func (d *DB) NamedQueryRowContext(ctx context.Context, query string, args map[string]any) (row *Row) {
	query, posArgs := convertToPositionalParameters(query, args)
	return d.QueryRowContext(ctx, query, posArgs...)
}
