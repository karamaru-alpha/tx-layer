package main

import (
	"context"
	"os"

	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/layer-tx/context-pattern/infra/mysql"
	"github.com/karamaru-alpha/layer-tx/context-pattern/infra/mysql/repository"
	"github.com/karamaru-alpha/layer-tx/context-pattern/usecase"
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

	txManager := mysql.NewTransactionManager(mysqlDB)
	userRepository := repository.NewUserRepository(mysqlDB)
	userInteractor := usecase.NewUserInteractor(txManager, userRepository)

	// Scenario Test
	ctx := context.Background()
	userID := "user_id"
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
