package sql

import (
	"github.com/CherkashinEvgeny/gosql/internal"
)

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
