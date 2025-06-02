package router

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	boxController "github.com/minminseo/recall-setter/controller/box"
	categoryController "github.com/minminseo/recall-setter/controller/category"
	patternController "github.com/minminseo/recall-setter/controller/pattern"
	userController "github.com/minminseo/recall-setter/controller/user"
)

func NewRouter(uc userController.IUserController, cc categoryController.ICategoryController, bc boxController.IBoxController, pc patternController.IPatternController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//CORS設定：フロントエンドからのリクエストを許可する
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken}, // 許可するヘッダーを指定
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"}, // 許可をしたいメソッドを設定
		AllowCredentials: true,                                     // Cookieの送受信を可能にする
	}))

	//CSRF対策：CookieとTokenで不正リクエストを防ぐ
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		// CookieSameSite: http.SameSiteNoneMode, // Postmanの動作確認ができたらこのモード
		CookieSameSite: http.SameSiteDefaultMode, // Postmanで動作確認する時はこのモード（セキュア属性をFalseにする）
		// CookieSameSite: http.SameSiteNoneMode,
		//CookieMaxAge:   60,Cookieの有効期限を設定するならこれ
	}))

	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)

	u := e.Group("/user")
	u.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
		ContextKey:  "user",
	}))

	u.GET("", uc.GetUserSetting)
	u.PUT("", uc.UpdateSetting)
	u.PUT("/password", uc.UpdatePassword)

	catGroup := e.Group("/categories")
	catGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
		ContextKey:  "user",
	}))

	catGroup.POST("", cc.CreateCategory)
	catGroup.GET("", cc.GetCategories)
	catGroup.PUT("/:id", cc.UpdateCategory)
	catGroup.DELETE("/:id", cc.DeleteCategory)

	boxGroup := e.Group("/categories/:category_id/boxes")
	boxGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
		ContextKey:  "user",
	}))
	boxGroup.POST("", bc.CreateBox)
	boxGroup.GET("", bc.GetBoxes)
	boxGroup.PUT("/:id", bc.UpdateBox)
	boxGroup.DELETE("/:id", bc.DeleteBox)

	patternGroup := e.Group("/patterns")
	patternGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))

	patternGroup.POST("", pc.CreatePattern)
	patternGroup.GET("", pc.GetPatterns)
	patternGroup.PUT("/:id", pc.UpdatePattern)
	patternGroup.DELETE("/:id", pc.DeletePattern)

	return e

}
