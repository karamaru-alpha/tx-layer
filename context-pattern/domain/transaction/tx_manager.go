package transaction

import "context"

type TxManager interface {
	Transaction(context.Context, func(context.Context) error) error
}
