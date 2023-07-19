package base

import (
	"context"
)

func Value(value any) (option Option) {
	return valueOption{value}
}

type valueOption struct {
	value any
}

func (o valueOption) Apply(factory Factory) (newFactory Factory) {
	return valueOptionFactory{factory, o.value}
}

type valueOptionFactory struct {
	next  Factory
	value any
}

func (f valueOptionFactory) Tx(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error) {
	options = append(options, f.value)
	return f.next.Tx(ctx, tx, options)
}
