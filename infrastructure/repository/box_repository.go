package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	boxDomain "github.com/minminseo/recall-setter/domain/box"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type boxRepository struct{}

func NewBoxRepository() boxDomain.IBoxRepository {
	return &boxRepository{}
}

func (r *boxRepository) Create(ctx context.Context, box *boxDomain.Box) error {
	q := db.GetQuery(ctx)

	parsedID, err := uuid.Parse(box.ID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedUserID, err := uuid.Parse(box.UserID)
	if err != nil {
		return err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	parsedCategoryID, err := uuid.Parse(box.CategoryID)
	if err != nil {
		return err
	}
	pgCategoryID := pgtype.UUID{Bytes: parsedCategoryID, Valid: true}

	parsedPatternID, err := uuid.Parse(box.PatternID)
	if err != nil {
		return err
	}
	pgPatternID := pgtype.UUID{Bytes: parsedPatternID, Valid: true}

	pgRegisteredAt := pgtype.Timestamptz{Time: box.RegisteredAt, Valid: true}
	pgEditedAt := pgtype.Timestamptz{Time: box.EditedAt, Valid: true}

	params := dbgen.CreateBoxParams{
		ID:           pgID,
		UserID:       pgUserID,
		CategoryID:   pgCategoryID,
		PatternID:    pgPatternID,
		Name:         box.Name,
		RegisteredAt: pgRegisteredAt,
		EditedAt:     pgEditedAt,
	}
	return q.CreateBox(ctx, params)
}

func (r *boxRepository) GetAllByCategoryID(ctx context.Context, categoryID, userID string) ([]*boxDomain.Box, error) {
	q := db.GetQuery(ctx)

	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}
	pgCategoryID := pgtype.UUID{Bytes: parsedCategoryID, Valid: true}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	rows, err := q.GetAllBoxesByCategoryID(ctx, dbgen.GetAllBoxesByCategoryIDParams{
		CategoryID: pgCategoryID,
		UserID:     pgUserID,
	})
	if err != nil {
		return nil, err
	}

	boxes := make([]*boxDomain.Box, len(rows))
	for i, row := range rows {
		id := uuid.UUID(row.ID.Bytes).String()
		uid := uuid.UUID(row.UserID.Bytes).String()
		cid := uuid.UUID(row.CategoryID.Bytes).String()
		pid := uuid.UUID(row.PatternID.Bytes).String()

		b, err := boxDomain.ReconstructBox(
			id,
			uid,
			cid,
			pid,
			row.Name,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
		boxes[i] = b
	}
	return boxes, nil
}

func (r *boxRepository) GetByID(ctx context.Context, boxID string, categoryID string, userID string) (*boxDomain.Box, error) {
	q := db.GetQuery(ctx)

	parsedID, err := uuid.Parse(boxID)
	if err != nil {
		return nil, err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}
	pgCategoryID := pgtype.UUID{Bytes: parsedCategoryID, Valid: true}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	params := dbgen.GetBoxByIDParams{
		ID:         pgID,
		CategoryID: pgCategoryID,
		UserID:     pgUserID,
	}
	row, err := q.GetBoxByID(ctx, params)
	if err != nil {
		return nil, err
	}
	b, err := boxDomain.ReconstructBox(
		boxID,
		userID,
		categoryID,
		uuid.UUID(row.PatternID.Bytes).String(),
		row.Name,
		row.RegisteredAt.Time,
		row.EditedAt.Time,
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// メモ：復習物ボックス内に復習物が存在しない場合のみにPatternIDの変更を許可する形式にする
func (r *boxRepository) Update(ctx context.Context, box *boxDomain.Box) error {
	q := db.GetQuery(ctx)

	parsedID, err := uuid.Parse(box.ID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedUserID, err := uuid.Parse(box.UserID)
	if err != nil {
		return err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	pgEditedAt := pgtype.Timestamptz{Time: box.EditedAt, Valid: true}

	params := dbgen.UpdateBoxParams{
		Name:     box.Name,
		EditedAt: pgEditedAt,
		ID:       pgID,
		UserID:   pgUserID,
	}
	return q.UpdateBox(ctx, params)
}

func (r *boxRepository) UpdateWithPatternID(ctx context.Context, box *boxDomain.Box) (int64, error) {
	q := db.GetQuery(ctx)

	parsedID, err := uuid.Parse(box.ID)
	if err != nil {
		return 0, err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedCategoryID, err := uuid.Parse(box.CategoryID)
	if err != nil {
		return 0, err
	}
	pgCategoryID := pgtype.UUID{Bytes: parsedCategoryID, Valid: true}

	parsedUserID, err := uuid.Parse(box.UserID)
	if err != nil {
		return 0, err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	parsedPatternID, err := uuid.Parse(box.PatternID)
	if err != nil {
		return 0, err
	}
	pgPatternID := pgtype.UUID{Bytes: parsedPatternID, Valid: true}

	pgEditedAt := pgtype.Timestamptz{Time: box.EditedAt, Valid: true}

	params := dbgen.UpdateBoxIfNoReviewItemsParams{
		PatternID:  pgPatternID,
		Name:       box.Name,
		EditedAt:   pgEditedAt,
		ID:         pgID,
		CategoryID: pgCategoryID,
		UserID:     pgUserID,
		BoxID:      pgID,
	}
	return q.UpdateBoxIfNoReviewItems(ctx, params)
}

func (r *boxRepository) Delete(ctx context.Context, boxID string, categoryID string, userID string) error {
	q := db.GetQuery(ctx)

	parsedID, err := uuid.Parse(boxID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	params := dbgen.DeleteBoxParams{
		ID:     pgID,
		UserID: pgUserID,
	}
	return q.DeleteBox(ctx, params)
}

func (r *boxRepository) GetBoxNamesByBoxIDs(ctx context.Context, ids []string) ([]*boxDomain.BoxName, error) {
	q := db.GetQuery(ctx)
	pgIDs := make([]pgtype.UUID, len(ids))
	for i, sid := range ids {
		u, err := uuid.Parse(sid)
		if err != nil {
			return nil, err
		}
		pgIDs[i] = pgtype.UUID{Bytes: u, Valid: true}
	}
	rows, err := q.GetBoxNamesByBoxIDs(ctx, pgIDs)
	if err != nil {
		return nil, err
	}
	out := make([]*boxDomain.BoxName, len(rows))
	for i, row := range rows {
		out[i] = &boxDomain.BoxName{
			BoxID:     uuid.UUID(row.ID.Bytes).String(),
			Name:      row.Name,
			PatternID: row.PatternID.String(),
		}
	}
	return out, nil
}
