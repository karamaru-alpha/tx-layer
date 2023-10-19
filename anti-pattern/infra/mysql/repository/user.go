package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/karamaru-alpha/layer-tx/anti-pattern/domain/entity"
	"github.com/karamaru-alpha/layer-tx/anti-pattern/domain/repository"
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

func (r *userRepository) LoadByPK(ctx context.Context, tx *sqlx.Tx, userID string) (*entity.User, error) {
	var user User
	if err := tx.GetContext(ctx, &user, "SELECT * FROM users WHERE user_id = ?", userID); err != nil {
		return nil, err
	}
	return user.toEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, tx *sqlx.Tx, user *entity.User) error {
	if _, err := tx.ExecContext(ctx, "UPDATE users SET name = ? WHERE user_id = ?", user.Name, user.UserID); err != nil {
		return err
	}
	return nil
}
