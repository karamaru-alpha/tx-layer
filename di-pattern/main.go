package main

import (
	"context"
	"os"

	"golang.org/x/exp/slog"

	"github.com/google/uuid"

	"github.com/karamaru-alpha/layer-tx/di-pattern/infra/mysql"
	"github.com/karamaru-alpha/layer-tx/di-pattern/infra/mysql/repository"
	"github.com/karamaru-alpha/layer-tx/di-pattern/usecase"
)

func main() {
	// DI
	mysqlDB, err := mysql.NewDB(&mysql.Config{
		Addr:     os.Getenv("MYSQL_ADDR"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DB:       os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	txManager := mysql.NewTxManager(mysqlDB)
	userRepository := repository.NewUserRepository()
	userInteractor := usecase.NewUserInteractor(txManager, userRepository)

	// Scenario Test
	ctx := context.Background()
	userID := uuid.New().String()
	if err := userInteractor.Create(ctx, userID, "old name"); err != nil {
		slog.Error(err.Error())
		return
	}
	if _, err := userInteractor.GetUser(ctx, userID); err != nil {
		slog.Error(err.Error())
		return
	}
	if err := userInteractor.UpdateName(ctx, userID, "new name"); err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("success!!!")
}
