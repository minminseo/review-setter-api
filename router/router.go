package router

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	userController "github.com/minminseo/recall-setter/controller/user"
)

func NewRouter(uc userController.IUserController) *echo.Echo {
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

	u.GET("/", uc.GetUserSetting)
	u.PUT("/", uc.UpdateSetting)
	u.PUT("/password", uc.UpdatePassword)
	return e

}
