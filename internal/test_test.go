package internal

import (
	"context"
	"testing"
)

func BenchmarkOption(b *testing.B) {
	b.ReportAllocs()
	m := New("test", testDb{
		options: []Option{
			{"1", "d"},
			{"2", "d"},
			{"3", "d"},
		},
	})
	for i := 0; i < b.N; i++ {
		_ = m.Transactional(context.Background(), func(ctx context.Context) (err error) {
			return nil
		}, Option{"3", "o"}, Option{"4", "o"})
	}
}

type testDb struct {
	options []Option
}

func (t testDb) Executor() (executor any) {
	return nil
}

func (t testDb) Tx(ctx context.Context, tx Tx, options []Option) (newTx Tx, err error) {
	_ = FindOption("3", t.options, options)
	return nil, nil
}
