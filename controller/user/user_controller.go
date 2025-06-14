package user

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	userUsecase "github.com/minminseo/recall-setter/usecase/user"
)

type userController struct {
	uu userUsecase.IUserUsecase
}

func NewUserController(uu userUsecase.IUserUsecase) IUserController {
	return &userController{uu: uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var request sighUpUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	input := userUsecase.CreateUserInput{
		Email:      request.Email,
		Password:   request.Password,
		Timezone:   request.Timezone,
		ThemeColor: request.ThemeColor,
		Language:   request.Language,
	}

	userRes, err := uc.uu.SignUp(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) VerifyEmail(c echo.Context) error {
	ctx := c.Request().Context()
	var request verifyEmailRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	input := userUsecase.VerifyEmailInput{
		Email: request.Email,
		Code:  request.Code,
	}

	// 認証に成功すると、Usecaseからログインレスポンスが返ってくる
	loginRes, err := uc.uu.VerifyEmail(ctx, input)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	// ログイン成功時と同様にCookieを設定
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = loginRes.Token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"theme_color": loginRes.ThemeColor,
		"language":    loginRes.Language,
	})
}

func (uc *userController) LogIn(c echo.Context) error {
	ctx := c.Request().Context()

	var request logInUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	input := userUsecase.LoginUserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	userRes, err := uc.uu.LogIn(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = userRes.Token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true // Postmanで動作確認する時はFalseにする
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"theme_color": userRes.ThemeColor,
		"language":    userRes.Language,
	})
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true // Postmanで動作確認する時はFalseにする
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)        // echoのコンテキストの中で、"csrf"というキーワードでtokenを取得
	return c.JSON(http.StatusOK, echo.Map{ /*上でstring型に型アサーションしてから
		JSONでクライアントにcsrfトークンをレスポンスする*/
		"csrf_token": token,
	})
}

func (uc *userController) GetUserSetting(c echo.Context) error {
	ctx := c.Request().Context()

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"]
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, "User ID not found in token")
	}
	userRes, err := uc.uu.GetUserSetting(ctx, userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) UpdateSetting(c echo.Context) error {
	ctx := c.Request().Context()
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	rawID, ok := claims["user_id"]
	if !ok || rawID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User ID not found in token"})
	}
	userID, ok := rawID.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid user ID in token"})
	}

	var request updateUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	input := userUsecase.UpdateUserInput{
		ID:         userID,
		Email:      request.Email,
		Timezone:   request.Timezone,
		ThemeColor: request.ThemeColor,
		Language:   request.Language,
	}

	userRes, err := uc.uu.UpdateSetting(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) UpdatePassword(c echo.Context) error {
	ctx := c.Request().Context()
	var request updatePasswordRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"]
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, "User ID not found in token")
	}
	request.ID = userID.(string)

	err := uc.uu.UpdatePassword(ctx, request.ID, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
