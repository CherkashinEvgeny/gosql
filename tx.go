package sqlx

import (
	"database/sql"
)

type Tx sql.Tx

func (t *Tx) Commit() (err error) {
	return (*sql.Tx)(t).Commit()
}

func (t *Tx) Rollback() (err error) {
	return (*sql.Tx)(t).Rollback()
}
