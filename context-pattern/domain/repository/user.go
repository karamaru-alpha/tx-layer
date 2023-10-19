package repository

import (
	"context"

	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/entity"
)

type UserRepository interface {
	LoadByPK(ctx context.Context, userID string) (*entity.User, error)
	Update(ctx context.Context, e *entity.User) error
	Insert(ctx context.Context, e *entity.User) error
}
