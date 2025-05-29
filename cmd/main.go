package main

import (
	"context"
	"log"
	"time"

	userController "github.com/minminseo/recall-setter/controller/user"
	userUsecase "github.com/minminseo/recall-setter/usecase/user"

	categoryController "github.com/minminseo/recall-setter/controller/category"
	categoryUsecase "github.com/minminseo/recall-setter/usecase/category"

	boxController "github.com/minminseo/recall-setter/controller/box"
	boxUsecase "github.com/minminseo/recall-setter/usecase/box"

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

	categoryRepository := repository.NewCategoryRepository(pool)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	categoryController := categoryController.NewCategoryController(categoryUsecase)

	boxRepository := repository.NewBoxRepository(pool)
	boxUsecase := boxUsecase.NewBoxUsecase(boxRepository)
	boxController := boxController.NewBoxController(boxUsecase)

	e := router.NewRouter(userController, categoryController, boxController)
	e.Logger.Fatal(e.Start(":8080"))
}
