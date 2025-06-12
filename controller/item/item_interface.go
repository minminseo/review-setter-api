package item

import "github.com/labstack/echo/v4"

type IItemController interface {
	CreateItem(c echo.Context) error
	UpdateItem(c echo.Context) error
	UpdateReviewDates(c echo.Context) error
	UpdateItemAsFinishedForce(c echo.Context) error
	UpdateReviewDateAsCompleted(c echo.Context) error
	UpdateReviewDateAsInCompleted(c echo.Context) error
	UpdateItemAsUnFinishedForce(c echo.Context) error
	DeleteItem(c echo.Context) error

	GetAllUnFinishedItemsByBoxID(c echo.Context) error
	GetAllUnFinishedUnclassifiedItemsByUserID(c echo.Context) error
	GetAllUnFinishedUnclassifiedItemsByCategoryID(c echo.Context) error

	CountItemsGroupedByBoxByUserID(c echo.Context) error
	CountUnclassifiedItemsGroupedByCategoryByUserID(c echo.Context) error
	CountUnclassifiedItemsByUserID(c echo.Context) error

	CountDailyDatesGroupedByBoxByUserID(c echo.Context) error
	CountDailyDatesUnclassifiedGroupedByCategoryByUserID(c echo.Context) error
	CountDailyDatesUnclassifiedByUserID(c echo.Context) error

	CountAllDailyReviewDates(c echo.Context) error

	GetAllDailyReviewDates(c echo.Context) error

	GetFinishedItemsByBoxID(c echo.Context) error
	GetUnclassfiedFinishedItemsByCategoryID(c echo.Context) error
	GetUnclassfiedFinishedItemsByUserID(c echo.Context) error
}
