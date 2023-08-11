package sql

import (
	"database/sql"

	"github.com/CherkashinEvgeny/gosql/internal"
)

const isolationLevelKey = "isolation"

type IsolationLevel = sql.IsolationLevel

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
