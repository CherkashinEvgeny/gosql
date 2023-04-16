package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TxOptions = sql.TxOptions

type TxController interface {
	Transactional(
		ctx context.Context,
		f func(ctx context.Context) error,
		options *TxOptions,
	) (beginErr error, err error, commitErr error, rollbackErr error)
	Begin(ctx context.Context, options *TxOptions) (txCtx context.Context, err error)
	End(ctx context.Context, txErr error) (rollbackErr error, commitErr error)
	Commit(ctx context.Context) (err error)
	Rollback(ctx context.Context) (err error)
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
	var tx *sql.Tx
	tx, err = (*sql.DB)(d).BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, d.txKey(), (*Tx)(tx)), nil
}

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
	txValue := ctx.Value(d.txKey())
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

func (d *DB) txKey() (key string) {
	return fmt.Sprintf("tx-%p", d)
}
