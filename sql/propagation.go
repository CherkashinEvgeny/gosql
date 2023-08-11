package sql

import (
	"github.com/CherkashinEvgeny/gosql/internal"
)

const propagationKey = "propagation"

func WithPropagation(propagation Propagation) Option {
	return internal.Option{
		Key:   propagationKey,
		Value: propagation,
	}
}

func getPropagation(options internal.Options) (propagation Propagation) {
	propagation, _ = options.Value(propagationKey).(Propagation)
	return propagation
}

type Propagation int

const (
	Never Propagation = iota - 2
	Supports
	Required
	Nested
	Mandatory
)
