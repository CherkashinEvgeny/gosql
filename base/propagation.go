package base

import (
	"context"
	"errors"
)

var (
	Required  Option = requiredPropagation{}
	Supports  Option = supportsPropagation{}
	Mandatory Option = mandatoryPropagation{}
	Never     Option = neverPropagation{}
)

type requiredPropagation struct {
}

func (r requiredPropagation) Apply(factory Factory) (newFactory Factory) {
	return requiredPropagationFactory{factory}
}

type requiredPropagationFactory struct {
	parent Factory
}

func (f requiredPropagationFactory) Executor() (executor any) {
	return f.parent.Executor()
}

func (f requiredPropagationFactory) Tx(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error) {
	if tx == nil {
		return f.parent.Tx(ctx, tx, options)
	}
	return nopTx{tx.Executor()}, nil
}

type supportsPropagation struct {
}

func (r supportsPropagation) Apply(factory Factory) (newFactory Factory) {
	return supportsPropagationFactory{factory}
}

type supportsPropagationFactory struct {
	parent Factory
}

func (f supportsPropagationFactory) Executor() (executor any) {
	return f.parent.Executor()
}

func (f supportsPropagationFactory) Tx(_ context.Context, tx Tx, _ ...any) (newTx Tx, err error) {
	if tx != nil {
		return tx, nil
	}
	return nopTx{f.Executor()}, nil
}

type mandatoryPropagation struct {
}

func (r mandatoryPropagation) Apply(factory Factory) (newFactory Factory) {
	return mandatoryPropagationFactory{factory}
}

type mandatoryPropagationFactory struct {
	parent Factory
}

func (f mandatoryPropagationFactory) Executor() (executor any) {
	return f.parent.Executor()
}

func (f mandatoryPropagationFactory) Tx(_ context.Context, tx Tx, _ ...any) (newTx Tx, err error) {
	if tx == nil {
		return nil, transactionExpected
	}
	return nopTx{tx.Executor()}, nil
}

var transactionExpected = errors.New("transaction expected")

type neverPropagation struct {
}

func (r neverPropagation) Apply(factory Factory) (newFactory Factory) {
	return neverPropagationFactory{factory}
}

type neverPropagationFactory struct {
	parent Factory
}

func (f neverPropagationFactory) Executor() (executor any) {
	return f.parent.Executor()
}

func (f neverPropagationFactory) Tx(_ context.Context, tx Tx, _ ...any) (newTx Tx, err error) {
	if tx != nil {
		return nil, transactionIsNotExpected
	}
	return nil, nil
}

var transactionIsNotExpected = errors.New("transaction is not expected")

type nopTx struct {
	executor any
}

func (n nopTx) Executor() (executor any) {
	return n.executor
}

func (n nopTx) Commit(_ context.Context) (err error) {
	return nil
}

func (n nopTx) Rollback(_ context.Context) {
}
