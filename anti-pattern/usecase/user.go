package usecase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/layer-tx/anti-pattern/domain/repository"
)

type UserInteractor interface {
	UpdateName(ctx context.Context, userID, name string) error
}

type userInteractor struct {
	db             *sqlx.DB
	userRepository repository.UserRepository
}

func NewUserInteractor(
	db *sqlx.DB,
	userRepository repository.UserRepository,
) UserInteractor {
	return &userInteractor{
		db,
		userRepository,
	}
}

// UpdateName NOTE: usecase層でDB情報の関心を持ってしまっている
func (i *userInteractor) UpdateName(ctx context.Context, userID, name string) error {
	tx, err := i.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// panic -> rollback
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			panic(p)
		}
		// error -> rollback
		if err != nil {
			if e := tx.Rollback(); e != nil {
				slog.ErrorContext(ctx, "failed to MySQL Rollback")
			}
			return
		}
		// success -> commit
		if e := tx.Commit(); e != nil {
			slog.ErrorContext(ctx, "failed to MySQL Commit")
		}
	}()

	user, err := i.userRepository.SelectByPK(ctx, tx, userID)
	if err != nil {
		return err
	}
	user.Name = name
	if err := i.userRepository.Update(ctx, tx, user); err != nil {
		return err
	}

	return nil
}
