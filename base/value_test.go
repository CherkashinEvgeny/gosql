package base

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDb struct {
	handler func(ctx context.Context, tx Tx, options Valuer) (newTx Tx, err error)
}

func (t testDb) Executor() (executor any) {
	return "db"
}

func (t testDb) Tx(ctx context.Context, tx Tx, options Valuer) (newTx Tx, err error) {
	return t.handler(ctx, tx, options)
}

func TestValueGlobalPropagation(t *testing.T) {
	manager := New("test", testDb{
		handler: func(ctx context.Context, tx Tx, options Valuer) (newTx Tx, err error) {
			assert.Equal(t, "global", options.Value("global"))
			return tx, nil
		},
	}, []Option{Value("global", "global")}, []Option{})
	_ = manager.Transactional(context.Background(), func(ctx context.Context) (err error) {
		return nil
	})
}

func TestValue(t *testing.T) {
	var n any
	v, ok := n.(sql.IsolationLevel)
	fmt.Println(v, ok)
}

func TestValueLocalPropagation(t *testing.T) {
	manager := New("test", testDb{
		handler: func(ctx context.Context, tx Tx, options Valuer) (newTx Tx, err error) {
			assert.Equal(t, "local", options.Value("local"))
			return tx, nil
		},
	}, []Option{}, []Option{})
	_ = manager.Transactional(context.Background(), func(ctx context.Context) (err error) {
		return nil
	}, Value("local", "local"))
}

func TestValueGlobalAndLocalPropagation(t *testing.T) {
	manager := New("test", testDb{
		handler: func(ctx context.Context, tx Tx, options Valuer) (newTx Tx, err error) {
			assert.Equal(t, "local", options.Value("local"))
			assert.Equal(t, "global", options.Value("global"))
			return tx, nil
		},
	}, []Option{Value("local", "local")}, []Option{})
	_ = manager.Transactional(context.Background(), func(ctx context.Context) (err error) {
		return nil
	}, Value("local", "local"))
}
