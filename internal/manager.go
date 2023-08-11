package internal

import (
	"context"
	"fmt"
)

type Manager struct {
	key any
	db  Db
}

func New(key any, db Db) (manager Manager) {
	return Manager{
		key: key,
		db:  db,
	}
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) (err error), options ...Option) (err error) {
	tx := m.extractTxFromContext(ctx)
	tx, err = m.db.Tx(ctx, tx, options)
	if err != nil {
		return BeginError{err}
	}
	ctx = m.putTxToContext(ctx, tx)
	if tx == nil {
		return f(ctx)
	}
	return m.transactional(ctx, tx, f)
}

func (m Manager) transactional(ctx context.Context, tx Tx, f func(ctx context.Context) (err error)) (err error) {
	panicked := true
	defer func() {
		if panicked {
			_ = tx.Rollback(ctx)
		}
	}()
	err = f(ctx)
	panicked = false
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return RollbackError{rollbackErr, err}
		}
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return CommitError{err}
	}
	return nil
}

func (m Manager) Executor(ctx context.Context) (executor any) {
	tx := m.extractTxFromContext(ctx)
	if tx == nil {
		return m.db.Executor()
	}
	return tx.Executor()
}

func (m Manager) extractTxFromContext(ctx context.Context) (tx Tx) {
	contextAny := ctx.Value(m.key)
	if contextAny == nil {
		return nil
	}
	tx, _ = contextAny.(Tx)
	return tx
}

func (m Manager) putTxToContext(ctx context.Context, tx Tx) (newCtx context.Context) {
	return context.WithValue(ctx, m.key, tx)
}

type BeginError struct {
	cause error
}

func (e BeginError) Cause() error {
	return e.cause
}

func (e BeginError) Unwrap() error {
	return e.cause
}

func (e BeginError) Error() string {
	return fmt.Sprintf("begin: %v", e.cause)
}

type CommitError struct {
	cause error
}

func (e CommitError) Cause() error {
	return e.cause
}

func (e CommitError) Unwrap() error {
	return e.cause
}

func (e CommitError) Error() string {
	return fmt.Sprintf("commit: %v", e.cause)
}

type RollbackError struct {
	cause error
	tx    error
}

func (e RollbackError) Tx() error {
	return e.tx
}

func (e RollbackError) Cause() error {
	return e.cause
}

func (e RollbackError) Unwrap() error {
	return e.cause
}

func (e RollbackError) Error() string {
	return fmt.Sprintf("rollback: %v", e.cause)
}
