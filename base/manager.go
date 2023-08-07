package base

import (
	"context"
	"fmt"
)

type Executor interface {
	Executor() (executor any)
}

type Factory interface {
	Tx(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error)
}

type Db interface {
	Executor
	Factory
}

type Tx interface {
	Executor
	Commit(ctx context.Context) (err error)
	Rollback(ctx context.Context)
}

type Option interface {
	Apply(factory Factory) (newFactory Factory)
}

type Manager struct {
	key     any
	db      Db
	options []Option
}

func New(key any, db Db, options ...Option) (manager Manager) {
	return Manager{
		key:     key,
		db:      db,
		options: options,
	}
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) (err error), options ...Option) (err error) {
	var factory Factory = m.db
	for _, option := range options {
		factory = option.Apply(factory)
	}
	tx, err := factory.Tx(ctx, m.extractTxFromContext(ctx))
	if err != nil {
		return BeginError{err}
	}
	if tx != nil {
		ctx = m.putTxToContext(ctx, tx)
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback(ctx)
				panic(r)
			}
			if err != nil {
				tx.Rollback(ctx)
				return
			}
			err = tx.Commit(ctx)
			if err != nil {
				err = CommitError{err}
				return
			}
		}()
	}
	return f(ctx)
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
	err error
}

func (e CommitError) Cause() error {
	return e.err
}

func (e CommitError) Unwrap() error {
	return e.err
}

func (e CommitError) Error() string {
	return fmt.Sprintf("commit: %v", e.err)
}
