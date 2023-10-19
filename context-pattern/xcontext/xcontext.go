package xcontext

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	Tx *sqlx.Tx
}

type keyConstraint interface {
	Transaction
}

type key[T keyConstraint] struct{}

func WithValue[T keyConstraint](ctx context.Context, val T) context.Context {
	return context.WithValue(ctx, key[T]{}, val)
}

func Value[T keyConstraint](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(key[T]{}).(T)
	return val, ok
}
