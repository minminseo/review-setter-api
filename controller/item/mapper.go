package item

import itemUsecase "github.com/minminseo/recall-setter/usecase/item"

func mapToItemResponse(items []*itemUsecase.GetItemOutput) []ItemResponse {
	res := make([]ItemResponse, len(items))
	for i, item := range items {
		reviewDates := make([]ReviewDateResponse, len(item.ReviewDates))
		for j, rd := range item.ReviewDates {
			reviewDates[j] = ReviewDateResponse{
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
		res[i] = ItemResponse{
			ItemID:       item.ItemID,
			UserID:       item.UserID,
			CategoryID:   item.CategoryID,
			BoxID:        item.BoxID,
			PatternID:    item.PatternID,
			Name:         item.Name,
			Detail:       item.Detail,
			LearnedDate:  item.LearnedDate,
			IsFinished:   item.IsFinished,
			RegisteredAt: item.RegisteredAt,
			EditedAt:     item.EditedAt,
			ReviewDates:  reviewDates,
		}
	}
	return res
}
