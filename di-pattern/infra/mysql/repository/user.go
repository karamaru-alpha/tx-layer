package repository

import (
	"context"

	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/entity"
	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/repository"
	"github.com/karamaru-alpha/layer-tx/di-pattern/domain/transaction"
	"github.com/karamaru-alpha/layer-tx/di-pattern/infra/mysql"
)

type userRepository struct {
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

type User struct {
	UserID string `db:"user_id"`
	Name   string
}

func (u *User) toEntity() *entity.User {
	return &entity.User{
		UserID: u.UserID,
		Name:   u.Name,
	}
}

func (r *userRepository) LoadByPK(ctx context.Context, _tx transaction.ROTx, userID string) (*entity.User, error) {
	tx, err := mysql.ExtractROTx(_tx)
	if err != nil {
		return nil, err
	}

	var user User
	if err := tx.GetContext(ctx, &user, "SELECT * FROM users WHERE user_id = ?", userID); err != nil {
		return nil, err
	}
	return user.toEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, _tx transaction.RWTx, e *entity.User) error {
	tx, err := mysql.ExtractRWTx(_tx)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE users SET name = ? WHERE user_id = ?", e.Name, e.UserID); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Insert(ctx context.Context, _tx transaction.RWTx, e *entity.User) error {
	tx, err := mysql.ExtractRWTx(_tx)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO users VALUES(?, ?)", e.UserID, e.Name); err != nil {
		return err
	}
	return nil
}
