package base

import (
	"context"
	"fmt"
)

type Manager struct {
	factory Factory
	stor    Store
}

func New(factory Factory) (manager Manager) {
	return Manager{
		factory: factory,
		stor:    keyStore{key: fmt.Sprintf("%p", factory)},
	}
}

type keyStore struct {
	key any
}

func (s keyStore) Get(ctx context.Context) (core Tx) {
	contextAny := ctx.Value(s.key)
	if contextAny == nil {
		return nil
	}
	core, _ = contextAny.(Tx)
	return core
}

func (s keyStore) Set(ctx context.Context, core Tx) (newCtx context.Context) {
	return context.WithValue(ctx, s.key, core)
}

func (m Manager) Transactional(ctx context.Context, f func(ctx context.Context) (err error), options ...Option) (err error) {
	factory := m.factory
	for _, option := range options {
		factory = option.Apply(factory)
	}
	tx, err := factory.Tx(ctx, m.stor.Get(ctx))
	if err != nil {
		return BeginError{err}
	}
	if tx != nil {
		ctx = m.stor.Set(ctx, tx)
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

func (m Manager) Get(ctx context.Context) (executor any) {
	tx := m.stor.Get(ctx)
	if tx == nil {
		return m.factory.Executor()
	}
	return tx.Executor()
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
