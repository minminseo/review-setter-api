package item

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	itemDomain "github.com/minminseo/recall-setter/domain/item"
	itemUsecase "github.com/minminseo/recall-setter/usecase/item"
)

type itemController struct {
	iu itemUsecase.IItemUsecase
}

func NewItemController(iu itemUsecase.IItemUsecase) IItemController {
	return &itemController{iu: iu}
}

// JWTトークンからUserIDを抽出するヘルパー関数
func getUserIDFromContext(c echo.Context) (string, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return "", errors.New("invalid token context")
	}
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", errors.New("user_id not found in token")
	}
	return userID, nil
}

// 基本CRUD
func (ic *itemController) CreateItem(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}

	var req CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := itemUsecase.CreateItemInput{
		UserID:                   userID,
		CategoryID:               req.CategoryID,
		BoxID:                    req.BoxID,
		PatternID:                req.PatternID,
		Name:                     req.Name,
		Detail:                   req.Detail,
		LearnedDate:              req.LearnedDate,
		IsMarkOverdueAsCompleted: req.IsMarkOverdueAsCompleted,
		Today:                    req.Today,
	}

	out, err := ic.iu.CreateItem(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習物の作成に失敗しました: " + err.Error()})
	}
	reviewDates := make([]ReviewDateResponse, len(out.Reviewdates))
	for i, rd := range out.Reviewdates {
		reviewDates[i] = ReviewDateResponse{
			ReviewDateID:         rd.DateID,
			UserID:               rd.UserID,
			ItemID:               rd.ItemID,
			StepNumber:           rd.StepNumber,
			InitialScheduledDate: rd.InitialScheduledDate,
			ScheduledDate:        rd.ScheduledDate,
			IsCompleted:          rd.IsCompleted,
		}
	}

	res := ItemResponse{
		ItemID:       out.ItemID,
		UserID:       out.UserID,
		CategoryID:   out.CategoryID,
		BoxID:        out.BoxID,
		PatternID:    out.PatternID,
		Name:         out.Name,
		Detail:       out.Detail,
		LearnedDate:  out.LearnedDate,
		IsFinished:   out.IsCompleted,
		RegisteredAt: out.RegisteredAt,
		EditedAt:     out.EditedAt,
		ReviewDates:  reviewDates,
	}

	return c.JSON(http.StatusCreated, res)

}

func (ic *itemController) UpdateItem(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")

	var req UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := itemUsecase.UpdateItemInput{
		ItemID:                   itemID,
		UserID:                   userID,
		CategoryID:               req.CategoryID,
		BoxID:                    req.BoxID,
		PatternID:                req.PatternID,
		Name:                     req.Name,
		Detail:                   req.Detail,
		LearnedDate:              req.LearnedDate,
		IsMarkOverdueAsCompleted: req.IsMarkOverdueAsCompleted,
		Today:                    req.Today,
	}

	out, err := ic.iu.UpdateItem(ctx, input)
	if err != nil {
		if errors.Is(err, itemDomain.ErrNoDiff) || errors.Is(err, itemDomain.ErrHasCompletedReviewDate) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習物の更新に失敗しました: " + err.Error()})
	}
	reviewDates := make([]ReviewDateResponse, len(out.ReviewDates))
	for i, rd := range out.ReviewDates {
		reviewDates[i] = ReviewDateResponse{
			ReviewDateID:         rd.ReviewDateID,
			UserID:               rd.UserID,
			CategoryID:           rd.CategoryID,
			BoxID:                rd.BoxID,
			ItemID:               rd.ItemID,
			StepNumber:           rd.StepNumber,
			InitialScheduledDate: rd.InitialScheduledDate,
			ScheduledDate:        rd.ScheduledDate,
			IsCompleted:          rd.IsCompleted,
		}
	}

	res := ItemResponse{
		ItemID:      out.ItemID,
		UserID:      out.UserID,
		CategoryID:  out.CategoryID,
		BoxID:       out.BoxID,
		PatternID:   out.PatternID,
		Name:        out.Name,
		Detail:      out.Detail,
		LearnedDate: out.LearnedDate,
		IsFinished:  out.IsFinished,
		EditedAt:    out.EditedAt,
		ReviewDates: reviewDates,
	}

	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) DeleteItem(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")

	if err := ic.iu.DeleteItem(ctx, itemID, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習物の削除に失敗しました: " + err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// 更新系
func (ic *itemController) UpdateReviewDates(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")
	reviewDateID := c.Param("review_date_id")

	var req UpdateReviewDatesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	steps := make([]itemUsecase.PatternStepInReviewDate, len(req.PatternSteps))
	for i, s := range req.PatternSteps {
		steps[i] = itemUsecase.PatternStepInReviewDate{
			PatternStepID: s.PatternStepID, UserID: s.UserID, PatternID: s.PatternID,
			StepNumber: s.StepNumber, IntervalDays: s.IntervalDays,
		}
	}

	input := itemUsecase.UpdateBackReviewDateInput{
		ReviewDateID:             reviewDateID,
		UserID:                   userID,
		CategoryID:               req.CategoryID,
		BoxID:                    req.BoxID,
		ItemID:                   itemID,
		StepNumber:               req.StepNumber,
		InitialScheduledDate:     req.InitialScheduledDate,
		RequestScheduledDate:     req.RequestScheduledDate,
		IsMarkOverdueAsCompleted: req.IsMarkOverdueAsCompleted,
		Today:                    req.Today,
		LearnedDate:              req.LearnedDate,
		PatternStepsInReviewDate: steps,
	}

	out, err := ic.iu.UpdateReviewDates(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習日の更新に失敗しました: " + err.Error()})
	}
	reviewDates := make([]ReviewDateResponse, len(out.ReviewDates))
	for i, rd := range out.ReviewDates {
		reviewDates[i] = ReviewDateResponse{
			ReviewDateID:         rd.ReviewDateID,
			UserID:               rd.UserID,
			CategoryID:           rd.CategoryID,
			BoxID:                rd.BoxID,
			ItemID:               rd.ItemID,
			StepNumber:           rd.StepNumber,
			InitialScheduledDate: rd.InitialScheduledDate,
			ScheduledDate:        rd.ScheduledDate,
			IsCompleted:          rd.IsCompleted,
		}
	}

	res := UpdateBackReviewDateResponse{
		ItemID:      out.ItemID,
		UserID:      out.UserID,
		IsFinished:  out.IsFinished,
		EditedAt:    out.EditedAt,
		ReviewDates: reviewDates,
	}
	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) UpdateItemAsFinishedForce(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")

	input := itemUsecase.UpdateItemAsFinishedForceInput{ItemID: itemID, UserID: userID}
	out, err := ic.iu.UpdateItemAsFinishedForce(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習物の完了処理に失敗しました: " + err.Error()})
	}
	res := UpdateItemAsFinishedForceResponse{
		ItemID:     out.ItemID,
		UserID:     out.UserID,
		IsFinished: out.IsFinished,
		EditedAt:   out.EditedAt,
	}

	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) UpdateReviewDateAsCompleted(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")
	reviewDateID := c.Param("review_date_id")

	var req UpdateReviewDateAsCompletedRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := itemUsecase.UpdateReviewDateAsCompletedInput{
		ReviewDateID: reviewDateID,
		UserID:       userID,
		ItemID:       itemID,
		StepNumber:   req.StepNumber,
	}

	out, err := ic.iu.UpdateReviewDateAsCompleted(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習日の完了処理に失敗しました: " + err.Error()})
	}
	res := UpdateReviewDateAsCompletedResponse{
		ReviewDateID: out.ReviewDateID,
		UserID:       out.UserID,
		IsCompleted:  out.IsCompleted,
		IsFinished:   out.IsFinished,
		EditedAt:     out.EditedAt,
	}

	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) UpdateReviewDateAsInCompleted(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")
	reviewDateID := c.Param("review_date_id")

	var req UpdateReviewDateAsInCompletedRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := itemUsecase.UpdateReviewDateAsInCompletedInput{
		ReviewDateID: reviewDateID,
		UserID:       userID,
		ItemID:       itemID,
		StepNumber:   req.StepNumber,
	}

	out, err := ic.iu.UpdateReviewDateAsInCompleted(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習日の未完了処理に失敗しました: " + err.Error()})
	}
	res := UpdateReviewDateAsInCompletedResponse{
		ReviewDateID: out.ReviewDateID,
		UserID:       out.UserID,
		IsCompleted:  out.IsCompleted,
		IsFinished:   out.IsFinished,
		EditedAt:     out.EditedAt,
	}
	return c.JSON(http.StatusOK, res)
}

func (ic *itemController) UpdateItemAsUnFinishedForce(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	itemID := c.Param("item_id")

	var req UpdateItemAsUnFinishedForceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "リクエストの形式が正しくありません: " + err.Error()})
	}

	input := itemUsecase.UpdateItemAsUnFinishedForceInput{
		ItemID:      itemID,
		UserID:      userID,
		CategoryID:  req.CategoryID,
		BoxID:       req.BoxID,
		PatternID:   req.PatternID,
		LearnedDate: req.LearnedDate,
		Today:       req.Today,
	}

	out, err := ic.iu.UpdateItemAsUnFinishedForce(ctx, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習物の再開処理に失敗しました: " + err.Error()})
	}
	reviewDates := make([]ReviewDateResponse, len(out.ReviewDates))
	for i, rd := range out.ReviewDates {
		reviewDates[i] = ReviewDateResponse{
			ReviewDateID:         rd.ReviewDateID,
			UserID:               rd.UserID,
			CategoryID:           rd.CategoryID,
			BoxID:                rd.BoxID,
			ItemID:               rd.ItemID,
			StepNumber:           rd.StepNumber,
			InitialScheduledDate: rd.InitialScheduledDate,
			ScheduledDate:        rd.ScheduledDate,
			IsCompleted:          rd.IsCompleted,
		}
	}

	res := UpdateItemAsUnFinishedForceResponse{
		ItemID:      out.ItemID,
		UserID:      out.UserID,
		IsFinished:  out.IsFinished,
		EditedAt:    out.EditedAt,
		ReviewDates: reviewDates,
	}

	return c.JSON(http.StatusOK, res)

}

// 取得系
func (ic *itemController) GetAllUnFinishedItemsByBoxID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	boxID := c.Param("box_id")

	out, err := ic.iu.GetAllUnFinishedItemsByBoxID(ctx, boxID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックス内の復習物取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, mapToItemResponse(out))

}

func (ic *itemController) GetAllUnFinishedUnclassifiedItemsByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}

	out, err := ic.iu.GetAllUnFinishedUnclassifiedItemsByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "未分類の復習物取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, mapToItemResponse(out))

}

func (ic *itemController) GetAllUnFinishedUnclassifiedItemsByCategoryID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	categoryID := c.Param("category_id")

	out, err := ic.iu.GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx, userID, categoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリ内の未分類復習物取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, mapToItemResponse(out))

}

// カウント系
func (ic *itemController) CountItemsGroupedByBoxByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}

	out, err := ic.iu.CountItemsGroupedByBoxByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックス毎の復習物数取得に失敗しました: " + err.Error()})
	}
	res := make([]ItemCountGroupedByBoxResponse, len(out))
	for i, r := range out {
		res[i] = ItemCountGroupedByBoxResponse{
			CategoryID: r.CategoryID,
			BoxID:      r.BoxID,
			Count:      r.Count,
		}
	}
	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) CountUnclassifiedItemsGroupedByCategoryByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}

	out, err := ic.iu.CountUnclassifiedItemsGroupedByCategoryByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリ毎の未分類復習物数取得に失敗しました: " + err.Error()})
	}
	res := make([]UnclassifiedItemCountGroupedByCategoryResponse, len(out))
	for i, r := range out {
		res[i] = UnclassifiedItemCountGroupedByCategoryResponse{
			CategoryID: r.CategoryID,
			Count:      r.Count,
		}
	}
	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) CountUnclassifiedItemsByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}

	out, err := ic.iu.CountUnclassifiedItemsByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "未分類復習物数取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, CountResponse{Count: out})

}

func (ic *itemController) CountDailyDatesGroupedByBoxByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	today := c.QueryParam("today")

	out, err := ic.iu.CountDailyDatesGroupedByBoxByUserID(ctx, userID, today)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックス毎の今日の復習数取得に失敗しました: " + err.Error()})
	}
	res := make([]DailyCountGroupedByBoxResponse, len(out))
	for i, r := range out {
		res[i] = DailyCountGroupedByBoxResponse{
			CategoryID: r.CategoryID,
			BoxID:      r.BoxID,
			Count:      r.Count,
		}
	}
	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) CountDailyDatesUnclassifiedGroupedByCategoryByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	today := c.QueryParam("today")

	out, err := ic.iu.CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx, userID, today)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリ毎の今日の未分類復習数取得に失敗しました: " + err.Error()})
	}
	res := make([]UnclassifiedDailyDatesCountGroupedByCategoryResponse, len(out))
	for i, r := range out {
		res[i] = UnclassifiedDailyDatesCountGroupedByCategoryResponse{
			CategoryID: r.CategoryID,
			Count:      r.Count,
		}
	}
	return c.JSON(http.StatusOK, res)

}

func (ic *itemController) CountDailyDatesUnclassifiedByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	today := c.QueryParam("today")

	out, err := ic.iu.CountDailyDatesUnclassifiedByUserID(ctx, userID, today)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "今日の未分類復習数取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, CountResponse{Count: out})

}

// 今日の全復習日数を取得
func (ic *itemController) CountAllDailyReviewDates(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	today := c.QueryParam("today")

	count, err := ic.iu.CountAllDailyReviewDates(ctx, userID, today)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "今日の復習日数の取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, CountResponse{Count: count})

}

// 今日の復習日一覧取得
func (ic *itemController) GetAllDailyReviewDates(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	today := c.QueryParam("today")

	result, err := ic.iu.GetAllDailyReviewDates(ctx, userID, today)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "復習日の取得に失敗しました: " + err.Error()})
	}
	res := GetDailyReviewDatesResponse{}

	categories := make([]DailyReviewDatesGroupedByCategoryResponse, len(result.Categories))
	for i, cat := range result.Categories {
		boxes := make([]DailyReviewDatesGroupedByBoxResponse, len(cat.Boxes))
		for j, box := range cat.Boxes {
			reviewDates := make([]DailyReviewDatesByBoxResponse, len(box.ReviewDates))
			for k, rd := range box.ReviewDates {
				reviewDates[k] = DailyReviewDatesByBoxResponse{
					ReviewDateID:      rd.ReviewDateID,
					CategoryID:        rd.CategoryID,
					BoxID:             rd.BoxID,
					StepNumber:        rd.StepNumber,
					PrevScheduledDate: rd.PrevScheduledDate,
					ScheduledDate:     rd.ScheduledDate,
					NextScheduledDate: rd.NextScheduledDate,
					IsCompleted:       rd.IsCompleted,
					ItemName:          rd.ItemName,
					Detail:            rd.Detail,
					LearnedDate:       rd.LearnedDate,
					RegisteredAt:      rd.RegisteredAt,
					EditedAt:          rd.EditedAt,
				}
			}
			boxes[j] = DailyReviewDatesGroupedByBoxResponse{
				BoxID:        box.BoxID,
				CategoryID:   box.CategoryID,
				BoxName:      box.BoxName,
				ReviewDates:  reviewDates,
				TargetWeight: box.TargetWeight,
			}
		}

		unclassified := make([]UnclassifiedDailyReviewDatesGroupedByCategoryResponse, len(cat.UnclassifiedDailyReviewDatesByCategory))
		for j, rd := range cat.UnclassifiedDailyReviewDatesByCategory {
			unclassified[j] = UnclassifiedDailyReviewDatesGroupedByCategoryResponse{
				ReviewDateID:      rd.ReviewDateID,
				CategoryID:        rd.CategoryID,
				StepNumber:        rd.StepNumber,
				PrevScheduledDate: rd.PrevScheduledDate,
				ScheduledDate:     rd.ScheduledDate,
				NextScheduledDate: rd.NextScheduledDate,
				IsCompleted:       rd.IsCompleted,
				ItemName:          rd.ItemName,
				Detail:            rd.Detail,
				LearnedDate:       rd.LearnedDate,
				RegisteredAt:      rd.RegisteredAt,
				EditedAt:          rd.EditedAt,
			}
		}

		categories[i] = DailyReviewDatesGroupedByCategoryResponse{
			CategoryID:                             cat.CategoryID,
			CategoryName:                           cat.CategoryName,
			Boxes:                                  boxes,
			UnclassifiedDailyReviewDatesByCategory: unclassified,
		}
	}
	res.Categories = categories

	userUnclassified := make([]UnclassifiedDailyReviewDatesGroupedByUserResponse, len(result.DailyReviewDatesGroupedByUser))
	for i, rd := range result.DailyReviewDatesGroupedByUser {
		userUnclassified[i] = UnclassifiedDailyReviewDatesGroupedByUserResponse{
			ReviewDateID:      rd.ReviewDateID,
			StepNumber:        rd.StepNumber,
			PrevScheduledDate: rd.PrevScheduledDate,
			ScheduledDate:     rd.ScheduledDate,
			NextScheduledDate: rd.NextScheduledDate,
			IsCompleted:       rd.IsCompleted,
			ItemName:          rd.ItemName,
			Detail:            rd.Detail,
			LearnedDate:       rd.LearnedDate,
			RegisteredAt:      rd.RegisteredAt,
			EditedAt:          rd.EditedAt,
		}
	}
	res.DailyReviewDatesGroupedByUser = userUnclassified

	return c.JSON(http.StatusOK, res)
}

func (ic *itemController) GetFinishedItemsByBoxID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	boxID := c.Param("box_id")

	out, err := ic.iu.GetFinishedItemsByBoxID(ctx, boxID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "ボックス内の完了した復習物取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, mapToItemResponse(out))

}

func (ic *itemController) GetUnclassfiedFinishedItemsByCategoryID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}
	categoryID := c.Param("category_id")

	out, err := ic.iu.GetUnclassfiedFinishedItemsByCategoryID(ctx, userID, categoryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "カテゴリ内の未分類完了復習物取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, mapToItemResponse(out))

}

func (ic *itemController) GetUnclassfiedFinishedItemsByUserID(c echo.Context) error {
	ctx := c.Request().Context()
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "トークンにユーザーIDが含まれていません"})
	}

	out, err := ic.iu.GetUnclassfiedFinishedItemsByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "完了した復習物取得に失敗しました: " + err.Error()})
	}
	return c.JSON(http.StatusOK, mapToItemResponse(out))

}
