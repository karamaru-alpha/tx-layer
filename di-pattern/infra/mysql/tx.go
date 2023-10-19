package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/transaction"
)

type ROTx interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type RWTx interface {
	ROTx
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type rwTx struct {
	*sqlx.Tx
}

func (tx *rwTx) ROTxImpl() {}
func (tx *rwTx) RWTxImpl() {}

func ExtractRWTx(_tx transaction.RWTx) (RWTx, error) {
	tx, ok := _tx.(*rwTx)
	if !ok {
		return nil, errors.New("mysql RWTx is invalid")
	}
	return tx, nil
}

type roTx struct {
	// MysqlにはReadOnlyなTxオブジェクトが存在しない
	*sqlx.Tx
}

func (tx *roTx) ROTxImpl() {}

func ExtractROTx(_tx transaction.ROTx) (ROTx, error) {
	switch tx := _tx.(type) {
	case *roTx:
		return tx, nil
	case *rwTx: // ReadWriteTransaction内での呼び出しも許可する
		return tx, nil
	}
	return nil, errors.New("mysql ROTx is invalid")
}
