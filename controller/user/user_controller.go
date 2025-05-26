package user

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	userUsecase "github.com/minminseo/recall-setter/application/user"
)

type userController struct {
	uu userUsecase.IUserUsecase
}

func NewUserController(uu userUsecase.IUserUsecase) IUserController {
	return &userController{uu: uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	var request sighInUserRequest
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

	userRes, err := uc.uu.SignUp(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) LogIn(c echo.Context) error {
	var request logInUserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	input := userUsecase.LoginUserInput{
		Email:    request.Email,
		Password: request.Password,
	}

	userRes, err := uc.uu.LogIn(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = userRes.Token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = false // Postmanで動作確認する時はFalseにする
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
	cookie.Secure = false // Postmanで動作確認する時はFalseにする
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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"]
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, "User ID not found in token")
	}
	userRes, err := uc.uu.GetUserSetting(userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) UpdateSetting(c echo.Context) error {
	var request updateUserRequest
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

	input := userUsecase.UpdateUserInput{
		ID:         request.ID,
		Email:      request.Email,
		Timezone:   request.Timezone,
		ThemeColor: request.ThemeColor,
		Language:   request.Language,
	}

	userRes, err := uc.uu.UpdateSetting(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, userRes)
}

func (uc *userController) UpdatePassword(c echo.Context) error {
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

	err := uc.uu.UpdatePassword(request.ID, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
