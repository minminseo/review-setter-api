package main

import (
	"context"
	"log"
	"os"
	"time"

	itemDomain "github.com/minminseo/recall-setter/domain/item"
	userDomain "github.com/minminseo/recall-setter/domain/user"

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

	"github.com/minminseo/recall-setter/infrastructure/auth"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/mailer"
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

	encKey := os.Getenv("ENCRYPTION_KEY")
	if encKey == "" {
		log.Fatal("ENCRYPTION_KEYが設定されていません")
	}

	//暗号化と復号化のメソッドを持つcryptoServiceをインスタンス化
	cryptoService, err := userDomain.NewCryptoService(encKey)
	if err != nil {
		log.Fatalf("暗号化用の鍵の生成に失敗しました: %v", err)
	}

	hmacKey := os.Getenv("HMAC_SECRET_KEY")
	if hmacKey == "" {
		log.Fatal("HMAC_SECRET_KEYが設定されていません")
	}

	// 検索用のキー生成（ハッシュ値）のメソッドを持つhasherをインスタンス化
	hasher, err := userDomain.NewHasher(hmacKey)
	if err != nil {
		log.Fatalf("ハッシュ用のキーの生成に失敗しました: %v", err)
	}

	transactionManager := repository.NewTransactionManager(pool)

	// メール送信
	emailSender := mailer.NewSMTPEmailSender()

	// JWTトークン生成のためのサービス
	tokenGenerator := auth.NewJWTGenerator()

	// ドメインサービス
	scheduler := itemDomain.NewScheduler()

	// リポジトリ
	userRepository := repository.NewUserRepository()
	emailVerificationRepository := repository.NewEmailVerificationRepository()
	categoryRepository := repository.NewCategoryRepository()
	boxRepository := repository.NewBoxRepository()
	patternRepository := repository.NewPatternRepository()
	itemRepository := repository.NewItemRepository()

	// ユースケース
	userUsecase := userUsecase.NewUserUsecase(userRepository, emailVerificationRepository, transactionManager, cryptoService, hasher, emailSender, tokenGenerator)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	boxUsecase := boxUsecase.NewBoxUsecase(boxRepository)
	patternUsecase := patternUsecase.NewPatternUsecase(patternRepository, itemRepository, transactionManager)
	itemUsecase := itemUsecase.NewItemUsecase(categoryRepository, boxRepository, itemRepository, patternRepository, transactionManager, scheduler)

	// コントローラー
	userController := userController.NewUserController(userUsecase)
	categoryController := categoryController.NewCategoryController(categoryUsecase)
	boxController := boxController.NewBoxController(boxUsecase)
	patternController := patternController.NewPatternController(patternUsecase)
	itemController := itemController.NewItemController(itemUsecase)

	e := router.NewRouter(userController, categoryController, boxController, patternController, itemController)
	e.Logger.Fatal(e.Start(":8080"))
}
