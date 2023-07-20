package base

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDb struct {
}

func (t testDb) Executor() (executor any) {
	return "db"
}

type testFactory struct {
	handler func(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error)
}

func (t testFactory) Tx(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error) {
	return t.handler(ctx, tx, options...)
}

func TestValueGlobalPropagation(t *testing.T) {
	manager := New("test", testDb{}, testFactory{
		handler: func(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error) {
			assert.Equal(t, 1, len(options))
			assert.Equal(t, "global", options[0])
			return tx, nil
		},
	}, Value("global"))
	_ = manager.Transactional(context.Background(), func(ctx context.Context) (err error) {
		return nil
	})
}

func TestValueLocalPropagation(t *testing.T) {
	manager := New("test", testDb{}, testFactory{
		handler: func(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error) {
			assert.Equal(t, 1, len(options))
			assert.Equal(t, "local", options[0])
			return tx, nil
		},
	})
	_ = manager.Transactional(context.Background(), func(ctx context.Context) (err error) {
		return nil
	}, Value("local"))
}

func TestValueGlobalAndLocalPropagation(t *testing.T) {
	manager := New("test", testDb{}, testFactory{
		handler: func(ctx context.Context, tx Tx, options ...any) (newTx Tx, err error) {
			assert.Equal(t, 2, len(options))
			assert.Equal(t, "global", options[0])
			assert.Equal(t, "local", options[1])
			return tx, nil
		},
	}, Value("global"))
	_ = manager.Transactional(context.Background(), func(ctx context.Context) (err error) {
		return nil
	}, Value("local"))
}
