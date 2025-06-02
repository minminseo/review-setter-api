package box

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	boxUsecase "github.com/minminseo/recall-setter/usecase/box"
)

type boxController struct {
	bu boxUsecase.IBoxUsecase
}

func NewBoxController(bu boxUsecase.IBoxUsecase) IBoxController {
	return &boxController{bu: bu}
}

func (bc *boxController) CreateBox(c echo.Context) error {
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

	categoryIDParam := c.Param("category_id")

	var request CreateBoxRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := boxUsecase.CreateBoxInput{
		UserID:     userID,
		CategoryID: categoryIDParam,
		PatternID:  request.PatternID,
		Name:       request.Name,
	}
	boxRes, err := bc.bu.CreateBox(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックスの作成に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusCreated, boxRes)
}

func (bc *boxController) GetBoxes(c echo.Context) error {
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

	categoryIDParam := c.Param("category_id")
	boxesRes, err := bc.bu.GetBoxesByCategoryID(ctx, categoryIDParam, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックス一覧の取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, boxesRes)
}

func (bc *boxController) UpdateBox(c echo.Context) error {
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

	categoryIDParam := c.Param("category_id")
	boxID := c.Param("id")

	var req UpdateBoxRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := boxUsecase.UpdateBoxInput{
		ID:         boxID,
		UserID:     userID,
		CategoryID: categoryIDParam,
		PatternID:  req.PatternID,
		Name:       req.Name,
	}
	res, err := bc.bu.UpdateBox(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックスの更新に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (bc *boxController) DeleteBox(c echo.Context) error {
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

	categoryIDParam := c.Param("category_id")
	boxID := c.Param("id")

	err := bc.bu.DeleteBox(ctx, boxID, categoryIDParam, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックスの削除に失敗しました: " + err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
