package sql

import (
	"database/sql"

	"github.com/CherkashinEvgeny/gosql/base"
)

func Isolation(level sql.IsolationLevel) (option Option) {
	return base.Value(level)
}
