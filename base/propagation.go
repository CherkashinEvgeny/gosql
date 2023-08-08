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

type Propagation struct {
}

func (r Propagation) Name() (name string) {
	return "propagation"
}

func (r Propagation) Priority() (priority int) {
	return 0
}

type requiredPropagation struct {
	Propagation
}

func (r requiredPropagation) Apply(factory Factory) (newFactory Factory) {
	return requiredPropagationFactory{factory}
}

type requiredPropagationFactory struct {
	next Factory
}

func (f requiredPropagationFactory) Tx(ctx context.Context, tx Tx, options Valuer) (newTx Tx, err error) {
	if tx == nil {
		return f.next.Tx(ctx, tx, options)
	}
	return nopTx{tx}, nil
}

type supportsPropagation struct {
	Propagation
}

func (r supportsPropagation) Apply(factory Factory) (newFactory Factory) {
	return supportsPropagationFactory{factory}
}

type supportsPropagationFactory struct {
	next Factory
}

func (f supportsPropagationFactory) Tx(_ context.Context, tx Tx, _ Valuer) (newTx Tx, err error) {
	return tx, nil
}

type mandatoryPropagation struct {
	Propagation
}

func (r mandatoryPropagation) Apply(factory Factory) (newFactory Factory) {
	return mandatoryPropagationFactory{factory}
}

type mandatoryPropagationFactory struct {
	next Factory
}

func (f mandatoryPropagationFactory) Tx(_ context.Context, tx Tx, _ Valuer) (newTx Tx, err error) {
	if tx == nil {
		return nil, transactionExpected
	}
	return nopTx{tx}, nil
}

var transactionExpected = errors.New("transaction expected")

type neverPropagation struct {
	Propagation
}

func (r neverPropagation) Apply(factory Factory) (newFactory Factory) {
	return neverPropagationFactory{factory}
}

type neverPropagationFactory struct {
	next Factory
}

func (f neverPropagationFactory) Tx(_ context.Context, tx Tx, _ Valuer) (newTx Tx, err error) {
	if tx != nil {
		return nil, transactionIsNotExpected
	}
	return nil, nil
}

var transactionIsNotExpected = errors.New("transaction is not expected")

type nopTx struct {
	Tx
}

func (n nopTx) Commit(_ context.Context) (err error) {
	return nil
}

func (n nopTx) Rollback(_ context.Context) {
}
