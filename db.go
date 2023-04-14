package sqlx

import (
	"database/sql"
)

type DB sql.DB

func Open(driverName string, dataSource string) (db *DB, err error) {
	sqlDb, err := sql.Open(driverName, dataSource)
	return (*DB)(sqlDb), err
}
