package usecase

import (
	"context"

	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/entity"
	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/repository"
	"github.com/karamaru-alpha/layer-tx/context-pattern/domain/transaction"
)

type UserInteractor interface {
	Create(ctx context.Context, userID, name string) error
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	UpdateName(ctx context.Context, userID, name string) error
}

type userInteractor struct {
	txManager      transaction.TxManager
	userRepository repository.UserRepository
}

func NewUserInteractor(
	txManager transaction.TxManager,
	userRepository repository.UserRepository,
) UserInteractor {
	return &userInteractor{
		txManager,
		userRepository,
	}
}

func (i *userInteractor) Create(ctx context.Context, userID, name string) error {
	user := &entity.User{
		UserID: userID,
		Name:   name,
	}
	if err := i.userRepository.Insert(ctx, user); err != nil {
		return err
	}
	return nil
}

func (i *userInteractor) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	user, err := i.userRepository.LoadByPK(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *userInteractor) UpdateName(ctx context.Context, userID, name string) error {
	if err := i.txManager.Transaction(ctx, func(ctx context.Context) error {
		user, err := i.userRepository.LoadByPK(ctx, userID)
		if err != nil {
			return err
		}
		user.Name = name
		if err := i.userRepository.Update(ctx, user); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
