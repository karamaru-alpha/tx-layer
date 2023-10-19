package repository

import (
	"context"

	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/entity"
	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/transaction"
)

type UserRepository interface {
	LoadByPK(ctx context.Context, tx transaction.ROTx, userID string) (*entity.User, error)
	Update(ctx context.Context, tx transaction.RWTx, user *entity.User) error
}
