package category

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	categoryUsecase "github.com/minminseo/recall-setter/usecase/category"
)

type categoryController struct {
	cu categoryUsecase.ICategoryUsecase
}

func NewCategoryController(cu categoryUsecase.ICategoryUsecase) ICategoryController {
	return &categoryController{cu: cu}
}

func (cc *categoryController) CreateCategory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	rawID, ok := claims["user_id"]
	if !ok || rawID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	userID, ok := rawID.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークン内のユーザーIDが無効です"})
	}

	var request CreateCategoryRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := categoryUsecase.CreateCategoryInput{
		UserID: userID,
		Name:   request.Name,
	}

	categoryRes, err := cc.cu.CreateCategory(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの作成に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusCreated, categoryRes)
}

func (cc *categoryController) GetCategories(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	rawID, ok := claims["user_id"]
	if !ok || rawID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	userID, ok := rawID.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークン内のユーザーIDが無効です"})
	}

	categoriesRes, err := cc.cu.GetCategoriesByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, categoriesRes)
}

func (cc *categoryController) UpdateCategory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	rawID, ok := claims["user_id"]
	if !ok || rawID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	userID, ok := rawID.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークン内のユーザーIDが無効です"})
	}

	categoryIDParam := c.Param("id")
	if categoryIDParam == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "パスにカテゴリIDが必要です"})
	}

	var request UpdateCategoryRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := categoryUsecase.UpdateCategoryInput{
		ID:     categoryIDParam,
		UserID: userID,
		Name:   request.Name,
	}

	categoryRes, err := cc.cu.UpdateCategory(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの更新に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, categoryRes)
}

func (cc *categoryController) DeleteCategory(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	rawID, ok := claims["user_id"]
	if !ok || rawID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	userID, ok := rawID.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークン内のユーザーIDが無効です"})
	}
	categoryIDParam := c.Param("id")
	if categoryIDParam == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "パスにカテゴリIDが必要です"})
	}

	err := cc.cu.DeleteCategory(categoryIDParam, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの削除に失敗しました: " + err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
