package main

import (
	"context"
	"log"
	"time"

	userUsecase "github.com/minminseo/recall-setter/application/user"
	userController "github.com/minminseo/recall-setter/controller/user"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/repository"
	"github.com/minminseo/recall-setter/router"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, _, err := db.NewDB(ctx)
	if err != nil {
		log.Fatalf("DB接続に失敗しました: %v", err)
	}
	defer pool.Close()

	userRepository := repository.NewUserRepository(pool)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userController := userController.NewUserController(userUsecase)
	e := router.NewRouter(userController)
	e.Logger.Fatal(e.Start(":8080"))
}
