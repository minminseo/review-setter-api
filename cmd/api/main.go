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

	itemController "github.com/minminseo/recall-setter/controller/item"
	itemUsecase "github.com/minminseo/recall-setter/usecase/item"

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

	// リポジトリ
	userRepository := repository.NewUserRepository()
	emailVerificationRepository := repository.NewEmailVerificationRepository()
	categoryRepository := repository.NewCategoryRepository()
	boxRepository := repository.NewBoxRepository()
	patternRepository := repository.NewPatternRepository()
	itemRepository := repository.NewItemRepository()

	// ユースケース
	userUsecase := userUsecase.NewUserUsecase(userRepository, emailVerificationRepository, transactionManager)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	boxUsecase := boxUsecase.NewBoxUsecase(boxRepository)
	patternUsecase := patternUsecase.NewPatternUsecase(patternRepository, itemRepository, transactionManager)
	itemUsecase := itemUsecase.NewItemUsecase(categoryRepository, boxRepository, itemRepository, patternRepository, transactionManager)

	// コントローラー
	userController := userController.NewUserController(userUsecase)
	categoryController := categoryController.NewCategoryController(categoryUsecase)
	boxController := boxController.NewBoxController(boxUsecase)
	patternController := patternController.NewPatternController(patternUsecase)
	itemController := itemController.NewItemController(itemUsecase)

	e := router.NewRouter(userController, categoryController, boxController, patternController, itemController)
	e.Logger.Fatal(e.Start(":8080"))
}
