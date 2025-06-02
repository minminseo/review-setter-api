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

	patternController "github.com/minminseo/recall-setter/controller/pattern"
	patternUsecase "github.com/minminseo/recall-setter/usecase/pattern"

	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/repository"
	"github.com/minminseo/recall-setter/router"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewDB(ctx)
	if err != nil {
		log.Fatalf("DB接続に失敗しました: %v", err)
	}
	defer pool.Close()

	transactionManager := repository.NewTransactionManager(pool)

	userRepository := repository.NewUserRepository()
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userController := userController.NewUserController(userUsecase)

	categoryRepository := repository.NewCategoryRepository()
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	categoryController := categoryController.NewCategoryController(categoryUsecase)

	boxRepository := repository.NewBoxRepository()
	boxUsecase := boxUsecase.NewBoxUsecase(boxRepository)
	boxController := boxController.NewBoxController(boxUsecase)

	patternRepository := repository.NewPatternRepository()
	patternUsecase := patternUsecase.NewPatternUsecase(patternRepository, transactionManager)
	patternController := patternController.NewPatternController(patternUsecase)

	e := router.NewRouter(userController, categoryController, boxController, patternController)
	e.Logger.Fatal(e.Start(":8080"))
}
