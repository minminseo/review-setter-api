package pattern

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	patternUsecase "github.com/minminseo/recall-setter/usecase/pattern"
)

type patternController struct {
	pu patternUsecase.IPatternUsecase
}

func NewPatternController(pu patternUsecase.IPatternUsecase) IPatternController {
	return &patternController{pu: pu}
}

func (pc *patternController) CreatePattern(c echo.Context) error {
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

	var req CreatePatternRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	steps := make([]patternUsecase.CreatePatternStepInput, len(req.Steps))
	for i, s := range req.Steps {
		steps[i] = patternUsecase.CreatePatternStepInput{
			StepNumber:   s.StepNumber,
			IntervalDays: s.IntervalDays,
		}
	}
	input := patternUsecase.CreatePatternInput{
		UserID:       userID,
		Name:         req.Name,
		TargetWeight: req.TargetWeight,
		Steps:        steps,
	}

	out, err := pc.pu.CreatePattern(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "パターンの作成に失敗しました: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, out)
}

func (pc *patternController) GetPatterns(c echo.Context) error {
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

	results, err := pc.pu.GetPatternsByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "パターンの取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

func (pc *patternController) UpdatePattern(c echo.Context) error {
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

	patternID := c.Param("id")
	if patternID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "パスにパターンIDが必要です"})
	}

	var req UpdatePatternRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	steps := make([]patternUsecase.UpdatePatternStepInput, len(req.Steps))
	for i, s := range req.Steps {
		steps[i] = patternUsecase.UpdatePatternStepInput{
			StepID:       s.StepID,
			PatternID:    patternID,
			StepNumber:   s.StepNumber,
			IntervalDays: s.IntervalDays,
		}
	}
	input := patternUsecase.UpdatePatternInput{
		PatternID:    patternID,
		UserID:       userID,
		Name:         req.Name,
		TargetWeight: req.TargetWeight,
		Steps:        steps,
	}

	out, err := pc.pu.UpdatePattern(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "パターンの更新に失敗しました: " + err.Error()})
	}

	return c.JSON(http.StatusOK, out)
}

func (pc *patternController) DeletePattern(c echo.Context) error {
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

	patternID := c.Param("id")
	if patternID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "パスにパターンIDが必要です"})
	}

	if err := pc.pu.DeletePattern(ctx, patternID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "パターンの削除に失敗しました: " + err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
