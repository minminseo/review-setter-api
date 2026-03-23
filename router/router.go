package router

import (
	"log/slog"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	boxController "github.com/minminseo/recall-setter/controller/box"
	categoryController "github.com/minminseo/recall-setter/controller/category"
	itemController "github.com/minminseo/recall-setter/controller/item"

	patternController "github.com/minminseo/recall-setter/controller/pattern"
	userController "github.com/minminseo/recall-setter/controller/user"
)

func NewRouter(
	uc userController.IUserController,
	cc categoryController.ICategoryController,
	bc boxController.IBoxController,
	pc patternController.IPatternController,
	ic itemController.IItemController,
) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RequestID())
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:     true,
		LogMethod:        true,
		LogURI:           true,
		LogStatus:        true,
		LogRemoteIP:      true,
		LogLatency:       true,
		LogUserAgent:     true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogError:         true,
		HandleError:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("request_id", v.RequestID),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("ip", v.RemoteIP),
				slog.String("user_agent", v.UserAgent),
				slog.Float64("latency_ms", float64(v.Latency.Nanoseconds())/1e6),
				slog.String("bytes_in", v.ContentLength),
				slog.Int64("bytes_out", v.ResponseSize),
			}
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST", attrs...)
			} else {
				attrs = append(attrs, slog.Any("err", v.Error))
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "REQUEST_ERROR", attrs...)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowCredentials: true,
	}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieSameSite: http.SameSiteNoneMode,
		//CookieMaxAge:   60,Cookieの有効期限を設定するならこれ
	}))

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.POST("/verify-email", uc.VerifyEmail)
	e.GET("/csrf", uc.CsrfToken)
	e.POST("/password-reset/request", uc.RequestPasswordReset)
	e.POST("/password-reset/reset", uc.ResetPassword)

	authMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
		ContextKey:  "user",
	})

	userGroup := e.Group("/user")
	userGroup.Use(authMiddleware)
	{
		userGroup.GET("", uc.GetUserSetting)
		userGroup.PUT("", uc.UpdateSetting)
		userGroup.PUT("/password", uc.UpdatePassword)
	}

	// カテゴリー系
	categoryGroup := e.Group("/categories")
	categoryGroup.Use(authMiddleware)
	{
		categoryGroup.POST("", cc.CreateCategory)
		categoryGroup.GET("", cc.GetCategories)
		categoryGroup.PUT("/:id", cc.UpdateCategory)
		categoryGroup.DELETE("/:id", cc.DeleteCategory)
	}

	// ボックス系
	boxGroup := e.Group("/:category_id/boxes")
	boxGroup.Use(authMiddleware)
	{
		boxGroup.POST("", bc.CreateBox)
		boxGroup.GET("", bc.GetBoxes)
		boxGroup.PUT("/:id", bc.UpdateBox)
		boxGroup.DELETE("/:id", bc.DeleteBox)
	}

	// 復習パターン系
	patternGroup := e.Group("/patterns")
	patternGroup.Use(authMiddleware)
	{
		patternGroup.POST("", pc.CreatePattern)
		patternGroup.GET("", pc.GetPatterns)
		patternGroup.PUT("/:id", pc.UpdatePattern)
		patternGroup.DELETE("/:id", pc.DeletePattern)
	}

	// 復習打つ形
	itemGroup := e.Group("/items")
	itemGroup.Use(authMiddleware)
	{
		// 復習物の作成
		itemGroup.POST("", ic.CreateItem)

		// 復習物一覧取得系
		itemGroup.GET("/unclassified", ic.GetAllUnFinishedUnclassifiedItemsByUserID)
		itemGroup.GET("/:box_id", ic.GetAllUnFinishedItemsByBoxID)
		itemGroup.GET("/unclassified/:category_id", ic.GetAllUnFinishedUnclassifiedItemsByCategoryID)
		itemGroup.GET("/today", ic.GetAllDailyReviewDates)

		// 完了済み復習物一覧取得系
		itemGroup.GET("/finished/unclassified", ic.GetUnclassfiedFinishedItemsByUserID)
		itemGroup.GET("/finished/:box_id", ic.GetFinishedItemsByBoxID)
		itemGroup.GET("/finished/unclassified/:category_id", ic.GetUnclassfiedFinishedItemsByCategoryID)

		// 特定復習物への操作
		itemDetailGroup := itemGroup.Group("/:item_id")
		{
			itemDetailGroup.PUT("", ic.UpdateItem)
			itemDetailGroup.DELETE("", ic.DeleteItem)
			itemDetailGroup.PATCH("/finish", ic.UpdateItemAsFinishedForce)
			itemDetailGroup.PATCH("/unfinish", ic.UpdateItemAsUnFinishedForce)

			// 特定復習物に属する復習日への操作
			reviewDateGroup := itemDetailGroup.Group("/review-dates/:review_date_id")
			{
				// 復習日とその後の日付を再計算
				reviewDateGroup.PUT("", ic.UpdateReviewDates)
				// 復習日の完了状態を変更
				reviewDateGroup.PATCH("/complete", ic.UpdateReviewDateAsCompleted)
				reviewDateGroup.PATCH("/incomplete", ic.UpdateReviewDateAsInCompleted)
			}
		}
	}

	// データ概要系
	summaryGroup := e.Group("/summary")
	summaryGroup.Use(authMiddleware)
	{
		// 復習物の数：各カテゴリーのボックスごとの復習物の数、カテゴリーごとの未分類ボックスの復習物の数、ホーム画面の未分類ボックスの復習物の数
		summaryGroup.GET("/items/count/by-box", ic.CountItemsGroupedByBoxByUserID)
		summaryGroup.GET("/items/count/unclassified/by-category", ic.CountUnclassifiedItemsGroupedByCategoryByUserID)
		summaryGroup.GET("/items/count/unclassified", ic.CountUnclassifiedItemsByUserID)

		// 今日の復習内容の数 (日付はクエリパラメータで指定）：各カテゴリーのボックスごとの今日の復習内容の数、カテゴリーごとの未分類ボックスの今日の復習内容の数、ホーム画面の未分類ボックスの今日の復習内容の数
		summaryGroup.GET("/daily-reviews/count/by-box", ic.CountDailyDatesGroupedByBoxByUserID)
		summaryGroup.GET("/daily-reviews/count/unclassified/by-category", ic.CountDailyDatesUnclassifiedGroupedByCategoryByUserID)
		summaryGroup.GET("/daily-reviews/count/unclassified", ic.CountDailyDatesUnclassifiedByUserID)

		// 今日の全復習日数を取得
		summaryGroup.GET("/daily-reviews/count", ic.CountAllDailyReviewDates)
	}

	return e

}
