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
	ctx := c.Request().Context()

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

	categoryRes, err := cc.cu.CreateCategory(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの作成に失敗しました: " + err.Error()})
	}
	res := CategoryResponse{
		ID:           categoryRes.ID,
		UserID:       categoryRes.UserID,
		Name:         categoryRes.Name,
		RegisteredAt: categoryRes.RegisteredAt,
		EditedAt:     categoryRes.EditedAt,
	}
	return c.JSON(http.StatusCreated, res)

}

func (cc *categoryController) GetCategories(c echo.Context) error {
	ctx := c.Request().Context()

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

	categoriesRes, err := cc.cu.GetCategoriesByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの取得に失敗しました: " + err.Error()})
	}

	var res []CategoryResponse
	for _, cat := range categoriesRes {
		res = append(res, CategoryResponse{
			ID:           cat.ID,
			UserID:       cat.UserID,
			Name:         cat.Name,
			RegisteredAt: cat.RegisteredAt,
			EditedAt:     cat.EditedAt,
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (cc *categoryController) UpdateCategory(c echo.Context) error {
	ctx := c.Request().Context()

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

	categoryRes, err := cc.cu.UpdateCategory(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの更新に失敗しました: " + err.Error()})
	}

	res := UpdateCategoryResponse{
		ID:       categoryRes.ID,
		UserID:   categoryRes.UserID,
		Name:     categoryRes.Name,
		EditedAt: categoryRes.EditedAt,
	}
	return c.JSON(http.StatusOK, res)

}

func (cc *categoryController) DeleteCategory(c echo.Context) error {
	ctx := c.Request().Context()

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

	err := cc.cu.DeleteCategory(ctx, categoryIDParam, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリの削除に失敗しました: " + err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
