package user

import "github.com/labstack/echo/v4"

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(c echo.Context) error
	GetUserSetting(c echo.Context) error
	UpdateSetting(c echo.Context) error
	UpdatePassword(c echo.Context) error
	VerifyEmail(c echo.Context) error
}
