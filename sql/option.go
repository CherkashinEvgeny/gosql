package sql

import (
	"database/sql"
	"github.com/CherkashinEvgeny/gosql/internal"
)

const propagationKey = "propagation"

func WithPropagation(propagation Propagation) Option {
	return Option{
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

const isolationLevelKey = "isolation"

func WithIsolationLevel(level sql.IsolationLevel) (option Option) {
	return internal.Option{
		Key:   isolationLevelKey,
		Value: level,
	}
}

func getIsolationLevel(options internal.Options) (level sql.IsolationLevel) {
	level, _ = options.Value(isolationLevelKey).(sql.IsolationLevel)
	return level
}

const readonlyKey = "readonly"

func WithReadonly(readonly bool) (option Option) {
	return internal.Option{
		Key:   readonlyKey,
		Value: readonly,
	}
}

func getReadonly(options internal.Options) (readonly bool) {
	readonly, _ = options.Value(readonlyKey).(bool)
	return readonly
}
