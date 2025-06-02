package pattern

import "github.com/labstack/echo/v4"

type IPatternController interface {
	CreatePattern(c echo.Context) error
	GetPatterns(c echo.Context) error
	UpdatePattern(c echo.Context) error
	DeletePattern(c echo.Context) error
}
