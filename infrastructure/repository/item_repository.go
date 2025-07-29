package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	itemDomain "github.com/minminseo/recall-setter/domain/item"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type itemRepository struct{}

func NewItemRepository() itemDomain.IItemRepository {
	return &itemRepository{}
}

// stringのポインタをpgtype.UUIDに変換するヘルパー関数。nilの場合は無効なUUIDを返す。
func toNullableUUID(s *string) (pgtype.UUID, error) {
	if s == nil {
		return pgtype.UUID{Valid: false}, nil
	}
	return toUUID(*s)
}

func (r *itemRepository) CreateItem(ctx context.Context, item *itemDomain.Item) error {
	q := db.GetQuery(ctx)

	pgID, err := toUUID(item.ItemID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(item.UserID)
	if err != nil {
		return err
	}
	pgCategoryID, err := toNullableUUID(item.CategoryID)
	if err != nil {
		return err
	}
	pgBoxID, err := toNullableUUID(item.BoxID)
	if err != nil {
		return err
	}
	pgPatternID, err := toNullableUUID(item.PatternID)
	if err != nil {
		return err
	}

	params := dbgen.CreateItemParams{
		ID:           pgID,
		UserID:       pgUserID,
		CategoryID:   pgCategoryID,
		BoxID:        pgBoxID,
		PatternID:    pgPatternID,
		Name:         item.Name,
		Detail:       pgtype.Text{String: item.Detail, Valid: true},
		LearnedDate:  pgtype.Date{Time: item.LearnedDate, Valid: true},
		IsFinished:   item.IsFinished,
		RegisteredAt: pgtype.Timestamptz{Time: item.RegisteredAt, Valid: true},
		EditedAt:     pgtype.Timestamptz{Time: item.EditedAt, Valid: true},
	}
	return q.CreateItem(ctx, params)
}

func (r *itemRepository) CreateReviewdates(ctx context.Context, reviewdates []*itemDomain.Reviewdate) (int64, error) {
	q := db.GetQuery(ctx)

	params := make([]dbgen.CreateReviewDatesParams, len(reviewdates))
	rows := make([][]interface{}, len(reviewdates))

	for i, rd := range reviewdates {
		pgID, err := toUUID(rd.ReviewdateID)
		if err != nil {
			return 0, err
		}
		pgUserID, err := toUUID(rd.UserID)
		if err != nil {
			return 0, err
		}
		pgCategoryID, err := toNullableUUID(rd.CategoryID)
		if err != nil {
			return 0, err
		}
		pgBoxID, err := toNullableUUID(rd.BoxID)
		if err != nil {
			return 0, err
		}
		pgItemID, err := toUUID(rd.ItemID)
		if err != nil {
			return 0, err
		}

		params[i] = dbgen.CreateReviewDatesParams{
			ID:                   pgID,
			UserID:               pgUserID,
			CategoryID:           pgCategoryID,
			BoxID:                pgBoxID,
			ItemID:               pgItemID,
			StepNumber:           int16(rd.StepNumber), // #nosec G115
			InitialScheduledDate: pgtype.Date{Time: rd.InitialScheduledDate, Valid: true},
			ScheduledDate:        pgtype.Date{Time: rd.ScheduledDate, Valid: true},
			IsCompleted:          rd.IsCompleted,
		}

		rows[i] = []interface{}{
			params[i].ID,
			params[i].UserID,
			params[i].CategoryID,
			params[i].BoxID,
			params[i].ItemID,
			params[i].StepNumber,
			params[i].InitialScheduledDate,
			params[i].ScheduledDate,
			params[i].IsCompleted,
		}
	}

	columns := []string{"id", "user_id", "category_id", "box_id", "item_id", "step_number", "initial_scheduled_date", "scheduled_date", "is_completed"}
	return q.CopyFrom(
		ctx,
		pgx.Identifier{"review_dates"},
		columns,
		pgx.CopyFromRows(rows),
	)
}

func (r *itemRepository) GetItemByID(ctx context.Context, itemID string, userID string) (*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetItemByIDParams{
		ID:     pgItemID,
		UserID: pgUserID,
	}
	row, err := q.GetItemByID(ctx, params)
	if err != nil {
		return nil, err
	}

	var categoryID, boxID, patternID *string
	if row.CategoryID.Valid {
		idStr := uuid.UUID(row.CategoryID.Bytes).String()
		categoryID = &idStr
	}
	if row.BoxID.Valid {
		idStr := uuid.UUID(row.BoxID.Bytes).String()
		boxID = &idStr
	}
	if row.PatternID.Valid {
		idStr := uuid.UUID(row.PatternID.Bytes).String()
		patternID = &idStr
	}

	return itemDomain.ReconstructItem(
		uuid.UUID(row.ID.Bytes).String(),
		uuid.UUID(row.UserID.Bytes).String(),
		categoryID,
		boxID,
		patternID,
		row.Name,
		row.Detail.String,
		row.LearnedDate.Time,
		row.IsFinished,
		row.RegisteredAt.Time,
		row.EditedAt.Time,
	)
}

func (r *itemRepository) HasCompletedReviewDateByItemID(ctx context.Context, itemID string, userID string) (bool, error) {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return false, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return false, err
	}

	params := dbgen.HasCompletedReviewDateByItemIDParams{
		ItemID: pgItemID,
		UserID: pgUserID,
	}
	return q.HasCompletedReviewDateByItemID(ctx, params)
}

func (r *itemRepository) GetReviewDateIDsByItemID(ctx context.Context, itemID string, userID string) ([]string, error) {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}

	params := dbgen.GetReviewDateIDsByItemIDParams{
		ItemID: pgItemID,
		UserID: pgUserID,
	}
	ReviewDateID, err := q.GetReviewDateIDsByItemID(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]string, len(ReviewDateID))
	for i, row := range ReviewDateID {
		results[i] = uuid.UUID(row.Bytes).String()
	}
	return results, nil
}

func (r *itemRepository) UpdateItem(ctx context.Context, item *itemDomain.Item) error {
	q := db.GetQuery(ctx)
	pgID, err := toUUID(item.ItemID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(item.UserID)
	if err != nil {
		return err
	}
	pgCategoryID, err := toNullableUUID(item.CategoryID)
	if err != nil {
		return err
	}
	pgBoxID, err := toNullableUUID(item.BoxID)
	if err != nil {
		return err
	}
	pgPatternID, err := toNullableUUID(item.PatternID)
	if err != nil {
		return err
	}

	params := dbgen.UpdateItemParams{
		ID:          pgID,
		UserID:      pgUserID,
		CategoryID:  pgCategoryID,
		BoxID:       pgBoxID,
		PatternID:   pgPatternID,
		Name:        item.Name,
		Detail:      pgtype.Text{String: item.Detail, Valid: true},
		LearnedDate: pgtype.Date{Time: item.LearnedDate, Valid: true},
		IsFinished:  item.IsFinished,
		EditedAt:    pgtype.Timestamptz{Time: item.EditedAt, Valid: true},
	}
	return q.UpdateItem(ctx, params)
}

func (r *itemRepository) UpdateReviewDates(ctx context.Context, reviewdates []*itemDomain.Reviewdate, userID string) error {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}

	inputs := make([]string, len(reviewdates))
	for i, rd := range reviewdates {
		// UNNESTに渡すための(id,category_id,box_id,scheduled_date,is_completed)形式の文字列を生成
		var categoryID string
		if rd.CategoryID != nil {
			categoryID = *rd.CategoryID
		}
		var boxID string
		if rd.BoxID != nil {
			boxID = *rd.BoxID
		}
		inputs[i] = fmt.Sprintf("(%s,%s,%s,%s,%s,%t)", rd.ReviewdateID, categoryID, boxID, rd.InitialScheduledDate.Format("2006-01-02"), rd.ScheduledDate.Format("2006-01-02"), rd.IsCompleted)
	}

	params := dbgen.UpdateReviewDatesParams{
		UserID: pgUserID,
		Input:  inputs,
	}

	return q.UpdateReviewDates(ctx, params)
}

func (r *itemRepository) UpdateReviewDatesBack(ctx context.Context, reviewdates []*itemDomain.Reviewdate, userID string) error {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}

	inputs := make([]string, len(reviewdates))
	for i, rd := range reviewdates {
		// UNNESTに渡すための(id,category_id,box_id,scheduled_date,is_completed)形式の文字列を生成
		var categoryID string
		if rd.CategoryID != nil {
			categoryID = *rd.CategoryID
		}
		var boxID string
		if rd.BoxID != nil {
			boxID = *rd.BoxID
		}
		inputs[i] = fmt.Sprintf("(%s,%s,%s,%s,%s,%t)", rd.ReviewdateID, categoryID, boxID, rd.InitialScheduledDate.Format("2006-01-02"), rd.ScheduledDate.Format("2006-01-02"), rd.IsCompleted)
	}

	params := dbgen.UpdateReviewDatesBackParams{
		UserID: pgUserID,
		Input:  inputs,
	}

	return q.UpdateReviewDatesBack(ctx, params)
}

func (r *itemRepository) UpdateItemAsFinished(ctx context.Context, itemID string, userID string, editedAt time.Time) error {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	params := dbgen.UpdateItemAsFinishedParams{
		ID:       pgItemID,
		UserID:   pgUserID,
		EditedAt: pgtype.Timestamptz{Time: editedAt, Valid: true},
	}
	return q.UpdateItemAsFinished(ctx, params)
}

func (r *itemRepository) UpdateItemAsUnFinished(ctx context.Context, itemID string, userID string, editedAt time.Time) error {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	params := dbgen.UpdateItemAsUnfinishedParams{
		ID:       pgItemID,
		UserID:   pgUserID,
		EditedAt: pgtype.Timestamptz{Time: editedAt, Valid: true},
	}
	return q.UpdateItemAsUnfinished(ctx, params)
}

func (r *itemRepository) UpdateReviewDateAsCompleted(ctx context.Context, reviewdateID string, userID string) error {
	q := db.GetQuery(ctx)
	pgID, err := toUUID(reviewdateID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	params := dbgen.UpdateReviewDateAsCompletedParams{
		ID:     pgID,
		UserID: pgUserID,
	}
	return q.UpdateReviewDateAsCompleted(ctx, params)
}

func (r *itemRepository) UpdateReviewDateAsInCompleted(ctx context.Context, reviewdateID string, userID string) error {
	q := db.GetQuery(ctx)
	pgID, err := toUUID(reviewdateID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	params := dbgen.UpdateReviewDateAsInCompletedParams{
		ID:     pgID,
		UserID: pgUserID,
	}
	return q.UpdateReviewDateAsInCompleted(ctx, params)
}

func (r *itemRepository) GetReviewDatesByItemID(ctx context.Context, itemID string, userID string) ([]*itemDomain.Reviewdate, error) {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}

	params := dbgen.GetReviewDatesByItemIDParams{
		ItemID: pgItemID,
		UserID: pgUserID,
	}
	rows, err := q.GetReviewDatesByItemID(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]*itemDomain.Reviewdate, len(rows))
	for i, row := range rows {
		var categoryID, boxID *string
		if row.CategoryID.Valid {
			idStr := uuid.UUID(row.CategoryID.Bytes).String()
			categoryID = &idStr
		}
		if row.BoxID.Valid {
			idStr := uuid.UUID(row.BoxID.Bytes).String()
			boxID = &idStr
		}
		results[i], err = itemDomain.ReconstructReviewdate(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			categoryID,
			boxID,
			uuid.UUID(row.ItemID.Bytes).String(),
			int(row.StepNumber),
			row.InitialScheduledDate.Time,
			row.ScheduledDate.Time,
			row.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) DeleteItem(ctx context.Context, itemID string, userID string) error {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	params := dbgen.DeleteItemParams{
		ID:     pgItemID,
		UserID: pgUserID,
	}
	return q.DeleteItem(ctx, params)
}

func (r *itemRepository) DeleteReviewDates(ctx context.Context, itemID string, userID string) error {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return err
	}
	params := dbgen.DeleteReviewDatesParams{
		ItemID: pgItemID,
		UserID: pgUserID,
	}
	return q.DeleteReviewDates(ctx, params)
}

func (r *itemRepository) GetAllUnFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgBoxID, err := toUUID(boxID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetAllUnFinishedItemsByBoxIDParams{
		BoxID:  pgBoxID,
		UserID: pgUserID,
	}
	rows, err := q.GetAllUnFinishedItemsByBoxID(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]*itemDomain.Item, len(rows))
	for i, row := range rows {
		var categoryID, boxID, patternID *string
		if row.CategoryID.Valid {
			idStr := uuid.UUID(row.CategoryID.Bytes).String()
			categoryID = &idStr
		}
		if row.BoxID.Valid {
			idStr := uuid.UUID(row.BoxID.Bytes).String()
			boxID = &idStr
		}
		if row.PatternID.Valid {
			idStr := uuid.UUID(row.PatternID.Bytes).String()
			patternID = &idStr
		}
		results[i], err = itemDomain.ReconstructItem(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			categoryID,
			boxID,
			patternID,
			row.Name,
			row.Detail.String,
			row.LearnedDate.Time,
			row.IsFinished,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetAllReviewDatesByBoxID(ctx context.Context, boxID string, userID string) ([]*itemDomain.Reviewdate, error) {
	q := db.GetQuery(ctx)
	pgBoxID, err := toUUID(boxID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetAllReviewDatesByBoxIDParams{
		BoxID:  pgBoxID,
		UserID: pgUserID,
	}
	rows, err := q.GetAllReviewDatesByBoxID(ctx, params)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.Reviewdate, len(rows))
	for i, row := range rows {
		var categoryID, boxID *string
		if row.CategoryID.Valid {
			idStr := uuid.UUID(row.CategoryID.Bytes).String()
			categoryID = &idStr
		}
		if row.BoxID.Valid {
			idStr := uuid.UUID(row.BoxID.Bytes).String()
			boxID = &idStr
		}
		results[i], err = itemDomain.ReconstructReviewdate(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			categoryID,
			boxID,
			uuid.UUID(row.ItemID.Bytes).String(),
			int(row.StepNumber),
			row.InitialScheduledDate.Time,
			row.ScheduledDate.Time,
			row.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetAllUnFinishedUnclassifiedItemsByUserID(ctx context.Context, userID string) ([]*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := q.GetAllUnFinishedUnclassifiedItemsByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.Item, len(rows))
	for i, row := range rows {
		var patternID *string
		if row.PatternID.Valid {
			idStr := uuid.UUID(row.PatternID.Bytes).String()
			patternID = &idStr
		}
		results[i], err = itemDomain.ReconstructItem(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			nil, // Unclassified
			nil, // Unclassified
			patternID,
			row.Name,
			row.Detail.String,
			row.LearnedDate.Time,
			row.IsFinished,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetAllUnclassifiedReviewDatesByUserID(ctx context.Context, userID string) ([]*itemDomain.Reviewdate, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := q.GetAllUnclassifiedReviewDatesByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.Reviewdate, len(rows))
	for i, row := range rows {
		results[i], err = itemDomain.ReconstructReviewdate(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			nil, // Unclassified
			nil, // Unclassified
			uuid.UUID(row.ItemID.Bytes).String(),
			int(row.StepNumber),
			row.InitialScheduledDate.Time,
			row.ScheduledDate.Time,
			row.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx context.Context, categoryID string, userID string) ([]*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgCategoryID, err := toUUID(categoryID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetAllUnFinishedUnclassifiedItemsByCategoryIDParams{
		CategoryID: pgCategoryID,
		UserID:     pgUserID,
	}
	rows, err := q.GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]*itemDomain.Item, len(rows))
	for i, row := range rows {
		catIDStr := uuid.UUID(row.CategoryID.Bytes).String()
		var patternID *string
		if row.PatternID.Valid {
			idStr := uuid.UUID(row.PatternID.Bytes).String()
			patternID = &idStr
		}

		results[i], err = itemDomain.ReconstructItem(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			&catIDStr,
			nil, // Unclassified
			patternID,
			row.Name,
			row.Detail.String,
			row.LearnedDate.Time,
			row.IsFinished,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetAllUnclassifiedReviewDatesByCategoryID(ctx context.Context, categoryID string, userID string) ([]*itemDomain.Reviewdate, error) {
	q := db.GetQuery(ctx)
	pgCategoryID, err := toUUID(categoryID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetAllUnclassifiedReviewDatesByCategoryIDParams{
		CategoryID: pgCategoryID,
		UserID:     pgUserID,
	}
	rows, err := q.GetAllUnclassifiedReviewDatesByCategoryID(ctx, params)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.Reviewdate, len(rows))
	for i, row := range rows {
		catIDStr := uuid.UUID(row.CategoryID.Bytes).String()
		results[i], err = itemDomain.ReconstructReviewdate(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			&catIDStr,
			nil, // Unclassified
			uuid.UUID(row.ItemID.Bytes).String(),
			int(row.StepNumber),
			row.InitialScheduledDate.Time,
			row.ScheduledDate.Time,
			row.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) CountItemsGroupedByBoxByUserID(ctx context.Context, userID string) ([]*itemDomain.ItemCountGroupedByBox, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := q.CountItemsGroupedByBoxByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.ItemCountGroupedByBox, len(rows))
	for i, row := range rows {
		results[i] = &itemDomain.ItemCountGroupedByBox{
			CategoryID: uuid.UUID(row.CategoryID.Bytes).String(),
			BoxID:      uuid.UUID(row.BoxID.Bytes).String(),
			Count:      int(row.Count),
		}
	}
	return results, nil
}

func (r *itemRepository) CountUnclassifiedItemsGroupedByCategoryByUserID(ctx context.Context, userID string) ([]*itemDomain.UnclassifiedItemCountGroupedByCategory, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := q.CountUnclassifiedItemsGroupedByCategoryByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.UnclassifiedItemCountGroupedByCategory, len(rows))
	for i, row := range rows {
		results[i] = &itemDomain.UnclassifiedItemCountGroupedByCategory{
			CategoryID: uuid.UUID(row.CategoryID.Bytes).String(),
			Count:      int(row.Count),
		}
	}
	return results, nil
}

func (r *itemRepository) CountUnclassifiedItemsByUserID(ctx context.Context, userID string) (int, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return 0, err
	}
	counts, err := q.CountUnclassifiedItemsByUserID(ctx, pgUserID)
	if err != nil {
		return 0, err
	}
	if len(counts) == 0 {
		return 0, nil
	}
	return int(counts[0]), nil
}

func (r *itemRepository) CountDailyDatesGroupedByBoxByUserID(ctx context.Context, userID string, targetDate time.Time) ([]*itemDomain.DailyCountGroupedByBox, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.CountDailyDatesGroupedByBoxByUserIDParams{
		UserID:     pgUserID,
		TargetDate: pgtype.Date{Time: targetDate, Valid: true},
	}
	rows, err := q.CountDailyDatesGroupedByBoxByUserID(ctx, params)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.DailyCountGroupedByBox, len(rows))
	for i, row := range rows {
		results[i] = &itemDomain.DailyCountGroupedByBox{
			CategoryID: uuid.UUID(row.CategoryID.Bytes).String(),
			BoxID:      uuid.UUID(row.BoxID.Bytes).String(),
			Count:      int(row.Count),
		}
	}
	return results, nil
}

func (r *itemRepository) CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx context.Context, userID string, targetDate time.Time) ([]*itemDomain.UnclassifiedDailyDatesCountGroupedByCategory, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.CountDailyDatesUnclassifiedGroupedByCategoryByUserIDParams{
		UserID:     pgUserID,
		TargetDate: pgtype.Date{Time: targetDate, Valid: true},
	}
	rows, err := q.CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx, params)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.UnclassifiedDailyDatesCountGroupedByCategory, len(rows))
	for i, row := range rows {
		results[i] = &itemDomain.UnclassifiedDailyDatesCountGroupedByCategory{
			CategoryID: uuid.UUID(row.CategoryID.Bytes).String(),
			Count:      int(row.Count),
		}
	}
	return results, nil
}

func (r *itemRepository) CountDailyDatesUnclassifiedByUserID(ctx context.Context, userID string, targetDate time.Time) (int, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return 0, err
	}
	params := dbgen.CountDailyDatesUnclassifiedByUserIDParams{
		UserID:     pgUserID,
		TargetDate: pgtype.Date{Time: targetDate, Valid: true},
	}
	counts, err := q.CountDailyDatesUnclassifiedByUserID(ctx, params)
	if err != nil {
		return 0, err
	}
	if len(counts) == 0 {
		return 0, nil
	}
	return int(counts[0]), nil
}

func (r *itemRepository) IsPatternRelatedToItemByPatternID(ctx context.Context, patternID string, userID string) (bool, error) {
	q := db.GetQuery(ctx)
	pgPatternID, err := toUUID(patternID)
	if err != nil {
		return false, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return false, err
	}
	params := dbgen.IsPatternRelatedToItemByPatternIDParams{
		PatternID: pgPatternID,
		UserID:    pgUserID,
	}
	return q.IsPatternRelatedToItemByPatternID(ctx, params)
}

// EditedAtの取得専用
func (r *itemRepository) GetEditedAtByItemID(ctx context.Context, itemID string, userID string) (time.Time, error) {
	q := db.GetQuery(ctx)
	pgItemID, err := toUUID(itemID)
	if err != nil {
		return time.Time{}, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return time.Time{}, err
	}
	params := dbgen.GetEditedAtByItemIDParams{
		ID:     pgItemID,
		UserID: pgUserID,
	}
	editedAt, err := q.GetEditedAtByItemID(ctx, params)
	if err != nil {
		return time.Time{}, err
	}
	return editedAt.Time, nil
}

// 今日の全復習日数を取得
func (r *itemRepository) CountAllDailyReviewDates(ctx context.Context, userID string, targetDate time.Time) (int, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return 0, err
	}

	pgToday := pgtype.Date{Time: targetDate, Valid: true}

	params := dbgen.CountAllDailyReviewDatesParams{
		UserID:     pgUserID,
		TargetDate: pgToday,
	}

	counts, err := q.CountAllDailyReviewDates(ctx, params)
	if err != nil {
		return 0, err
	}
	return int(counts), nil
}

// 今日の復習日取得
func (r *itemRepository) GetAllDailyReviewDates(ctx context.Context, userID string, targetDate time.Time) ([]*itemDomain.DailyReviewDate, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}

	pgToday := pgtype.Date{Time: targetDate, Valid: true}

	params := dbgen.GetAllDailyReviewDatesParams{
		UserID: pgUserID,
		Today:  pgToday,
	}

	rows, err := q.GetAllDailyReviewDates(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]*itemDomain.DailyReviewDate, len(rows))
	for i, row := range rows {
		idStr := uuid.UUID(row.ID.Bytes).String()

		var categoryID *string
		if row.CategoryID.Valid {
			s := uuid.UUID(row.CategoryID.Bytes).String()
			categoryID = &s
		}

		var boxID *string
		if row.BoxID.Valid {
			s := uuid.UUID(row.BoxID.Bytes).String()
			boxID = &s
		}

		initialScheduledDate := row.InitialScheduledDate.Time

		var prev *time.Time
		if row.PrevScheduledDate.Valid {
			t := row.PrevScheduledDate.Time
			prev = &t
		}

		scheduled := row.ScheduledDate.Time

		var next *time.Time
		if row.NextScheduledDate.Valid {
			t := row.NextScheduledDate.Time
			next = &t
		}

		itemID := uuid.UUID(row.ItemID.Bytes).String()

		learnedDate := row.LearnedDate.Time

		detail := row.Detail.String

		results[i] = &itemDomain.DailyReviewDate{
			ReviewdateID:         idStr,
			CategoryID:           categoryID,
			BoxID:                boxID,
			StepNumber:           int(row.StepNumber),
			InitialScheduledDate: initialScheduledDate,
			PrevScheduledDate:    prev,
			ScheduledDate:        scheduled,
			NextScheduledDate:    next,
			IsCompleted:          row.IsCompleted,
			ItemID:               itemID,
			Name:                 row.Name,
			Detail:               detail,
			LearnedDate:          learnedDate,
			RegisteredAt:         row.RegisteredAt.Time,
			EditedAt:             row.EditedAt.Time,
		}
	}

	return results, nil
}

func (r *itemRepository) GetFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgBoxID, err := toUUID(boxID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetFinishedItemsByBoxIDParams{
		BoxID:  pgBoxID,
		UserID: pgUserID,
	}
	rows, err := q.GetFinishedItemsByBoxID(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]*itemDomain.Item, len(rows))
	for i, row := range rows {
		var categoryID, boxID, patternID *string
		if row.CategoryID.Valid {
			idStr := uuid.UUID(row.CategoryID.Bytes).String()
			categoryID = &idStr
		}
		if row.BoxID.Valid {
			idStr := uuid.UUID(row.BoxID.Bytes).String()
			boxID = &idStr
		}
		if row.PatternID.Valid {
			idStr := uuid.UUID(row.PatternID.Bytes).String()
			patternID = &idStr
		}
		results[i], err = itemDomain.ReconstructItem(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			categoryID,
			boxID,
			patternID,
			row.Name,
			row.Detail.String,
			row.LearnedDate.Time,
			row.IsFinished,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetUnclassfiedFinishedItemsByCategoryID(ctx context.Context, categoryID string, userID string) ([]*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgCategoryID, err := toUUID(categoryID)
	if err != nil {
		return nil, err
	}
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	params := dbgen.GetUnclassfiedFinishedItemsByCategoryIDParams{
		CategoryID: pgCategoryID,
		UserID:     pgUserID,
	}
	rows, err := q.GetUnclassfiedFinishedItemsByCategoryID(ctx, params)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.Item, len(rows))
	for i, row := range rows {
		catIDStr := uuid.UUID(row.CategoryID.Bytes).String()
		var patternID *string
		if row.PatternID.Valid {
			idStr := uuid.UUID(row.PatternID.Bytes).String()
			patternID = &idStr
		}
		results[i], err = itemDomain.ReconstructItem(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			&catIDStr,
			nil, // Unclassified
			patternID,
			row.Name,
			row.Detail.String,
			row.LearnedDate.Time,
			row.IsFinished,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (r *itemRepository) GetUnclassfiedFinishedItemsByUserID(ctx context.Context, userID string) ([]*itemDomain.Item, error) {
	q := db.GetQuery(ctx)
	pgUserID, err := toUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := q.GetUnclassfiedFinishedItemsByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}
	results := make([]*itemDomain.Item, len(rows))
	for i, row := range rows {
		var patternID *string
		if row.PatternID.Valid {
			idStr := uuid.UUID(row.PatternID.Bytes).String()
			patternID = &idStr
		}
		results[i], err = itemDomain.ReconstructItem(
			uuid.UUID(row.ID.Bytes).String(),
			uuid.UUID(row.UserID.Bytes).String(),
			nil, // Unclassified
			nil, // Unclassified
			patternID,
			row.Name,
			row.Detail.String,
			row.LearnedDate.Time,
			row.IsFinished,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
	}
	return results, nil
}
