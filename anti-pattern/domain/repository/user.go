package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/karamaru-alpha/layer-tx/anti-pattern/domain/entity"
)

// UserRepository NOTE: domain層でDB情報の関心を持ってしまっている
type UserRepository interface {
	SelectByPK(ctx context.Context, tx *sqlx.Tx, userID string) (*entity.User, error)
	Update(ctx context.Context, tx *sqlx.Tx, user *entity.User) error
}
