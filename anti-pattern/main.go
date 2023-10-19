package main

import (
	"context"
	"os"

	"golang.org/x/exp/slog"

	"github.com/karamaru-alpha/layer-tx/anti-pattern/infra/mysql"
	"github.com/karamaru-alpha/layer-tx/anti-pattern/infra/mysql/repository"
	"github.com/karamaru-alpha/layer-tx/anti-pattern/usecase"
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

	userRepository := repository.NewUserRepository()
	userInteractor := usecase.NewUserInteractor(mysqlDB, userRepository)

	// Scenario Test
	ctx := context.Background()
	userID := "user_id"
	if err := userInteractor.UpdateName(ctx, userID, "new name"); err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("success!!!")
}
