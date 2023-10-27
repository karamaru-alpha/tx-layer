package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/transaction"
)

type txManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) transaction.TxManager {
	return &txManager{db}
}

func (t *txManager) ReadWriteTransaction(ctx context.Context, f func(context.Context, transaction.RWTx) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// panic -> rollback
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			panic(p)
		}
		// error -> rollback
		if err != nil {
			if e := tx.Rollback(); e != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			return
		}
		// success -> commit
		if e := tx.Commit(); e != nil {
			slog.ErrorContext(ctx, "failed to MySQL Commit")
		}
	}()

	// ReadWriteTransactionオブジェクトを関数に渡す
	err = f(ctx, &rwTx{tx})
	if err != nil {
		return err
	}
	return nil
}

func (t *txManager) ReadOnlyTransaction(ctx context.Context, f func(context.Context, transaction.ROTx) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// panic -> rollback
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			panic(p)
		}
		// error -> rollback
		if err != nil {
			if e := tx.Rollback(); e != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			return
		}
		// success -> commit
		if e := tx.Commit(); e != nil {
			slog.ErrorContext(ctx, "failed to MySQL Commit")
		}
	}()

	// ReadOnlyTransactionオブジェクトを関数に渡す
	err = f(ctx, &roTx{tx})
	if err != nil {
		return err
	}
	return nil
}
