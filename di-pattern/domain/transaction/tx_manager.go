package transaction

import "context"

type ROTx interface {
	ROTxImpl()
}

type RWTx interface {
	ROTx
	RWTxImpl()
}

type TxManager interface {
	ReadOnlyTransaction(ctx context.Context, f func(ctx context.Context, tx ROTx) error) error
	ReadWriteTransaction(ctx context.Context, f func(ctx context.Context, tx RWTx) error) error
}
