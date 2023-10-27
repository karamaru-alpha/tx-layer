package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/entity"
	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/repository"
	"github.com/karamaru-alpha/layer-tx/context-pattern/infra/mysql"
	"github.com/karamaru-alpha/layer-tx/context-pattern/xcontext"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	return &userRepository{
		db,
	}
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

func (r *userRepository) SelectByPK(ctx context.Context, userID string) (*entity.User, error) {
	db := r.getMysqlDB(ctx)

	var user User
	if err := db.GetContext(ctx, &user, "SELECT * FROM users WHERE user_id = ?", userID); err != nil {
		return nil, err
	}
	return user.toEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, e *entity.User) error {
	db := r.getMysqlDB(ctx)

	if _, err := db.ExecContext(ctx, "UPDATE users SET name = ? WHERE user_id = ?", e.Name, e.UserID); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) getMysqlDB(ctx context.Context) mysql.DB {
	// contextにtxオブジェクトが存在すればそれを返却する
	if transaction, ok := xcontext.Value[xcontext.Transaction](ctx); ok {
		return transaction.Tx
	}
	// contextにtxオブジェクトが存在しなければDIされたdbを返却する
	return r.db
}
