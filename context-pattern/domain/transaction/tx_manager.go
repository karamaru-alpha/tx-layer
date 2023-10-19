package transaction

import "context"

type TxManager interface {
	Transaction(ctx context.Context, f func(context.Context) error) error
}
