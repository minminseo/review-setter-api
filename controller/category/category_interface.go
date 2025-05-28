package category

import "github.com/labstack/echo/v4"

type ICategoryController interface {
	CreateCategory(c echo.Context) error
	GetCategories(c echo.Context) error
	UpdateCategory(c echo.Context) error
	DeleteCategory(c echo.Context) error
}
