package sql

import (
	"context"
	"database/sql"

	"github.com/CherkashinEvgeny/gosql/internal"
)

type baseTx struct {
	tx      *sql.Tx
	options *sql.TxOptions
}

func (t baseTx) Executor() (executor any) {
	return t.tx
}

func (t baseTx) Parent() (db internal.Tx) {
	return nil
}

func (t baseTx) Commit(_ context.Context) (err error) {
	return t.tx.Commit()
}

func (t baseTx) Rollback(_ context.Context) (err error) {
	return t.tx.Rollback()
}

func nop(tx internal.Tx) (newTx internal.Nop) {
	return internal.Nop{Tx: tx}
}
