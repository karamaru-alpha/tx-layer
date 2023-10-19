package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/transaction"
	"github.com/karamaru-alpha/layer-tx/context-pattern/xcontext"
)

type txManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) transaction.TxManager {
	return &txManager{
		db,
	}
}

func (t *txManager) Transaction(ctx context.Context, f func(context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// panic
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			panic(p)
		}
		// error
		if err != nil {
			if e := tx.Rollback(); e != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			return
		}
		// success
		if e := tx.Commit(); e != nil {
			slog.ErrorContext(ctx, "failed to MySQL Commit")
		}
	}()

	// ContextにTxをセット
	ctx = xcontext.WithValue[xcontext.Transaction](ctx, xcontext.Transaction{
		Tx: tx,
	})
	err = f(ctx)
	if err != nil {
		return err
	}
	return nil
}
