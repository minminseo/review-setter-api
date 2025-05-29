package box

import "github.com/labstack/echo/v4"

type IBoxController interface {
	CreateBox(c echo.Context) error
	GetBoxes(c echo.Context) error
	UpdateBox(c echo.Context) error
	DeleteBox(c echo.Context) error
}
