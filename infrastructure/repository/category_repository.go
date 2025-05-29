package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	categoryDomain "github.com/minminseo/recall-setter/domain/category"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type categoryRepository struct {
	q *dbgen.Queries
}

func NewCategoryRepository(db dbgen.DBTX) categoryDomain.CategoryRepository {
	return &categoryRepository{
		q: dbgen.New(db),
	}
}

func (r *categoryRepository) Create(category *categoryDomain.Category) error {
	ctx := context.Background()

	parsedID, err := uuid.Parse(category.ID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedUserID, err := uuid.Parse(category.UserID)
	if err != nil {
		return err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	// 登録日時と更新日時をマッピング
	pgRegisteredAt := pgtype.Timestamptz{Time: category.RegisteredAt, Valid: true}
	pgEditedAt := pgtype.Timestamptz{Time: category.EditedAt, Valid: true}

	params := dbgen.CreateCategoryParams{
		ID:           pgID,
		UserID:       pgUserID,
		Name:         category.Name,
		RegisteredAt: pgRegisteredAt,
		EditedAt:     pgEditedAt,
	}

	return r.q.CreateCategory(ctx, params)
}

func (r *categoryRepository) GetAllByUserID(userID string) ([]*categoryDomain.Category, error) {
	ctx := context.Background()

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	rows, err := r.q.GetAllCategoriesByUserID(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	var categories []*categoryDomain.Category
	for _, row := range rows {
		catID := uuid.UUID(row.ID.Bytes).String()
		catUserID := uuid.UUID(row.UserID.Bytes).String()

		cat, err := categoryDomain.ReconstructCategory(
			catID,
			catUserID,
			row.Name,
			row.RegisteredAt.Time,
			row.EditedAt.Time,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}

func (r *categoryRepository) GetByID(categoryID string, userID string) (*categoryDomain.Category, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}
	params := dbgen.GetCategoryByIDParams{
		ID:     pgID,
		UserID: pgUserID,
	}
	row, err := r.q.GetCategoryByID(ctx, params)
	if err != nil {
		return nil, err
	}
	cd, err := categoryDomain.ReconstructCategory(
		categoryID,
		userID,
		row.Name,
		row.RegisteredAt.Time,
		row.EditedAt.Time,
	)
	if err != nil {
		return nil, err
	}
	return cd, nil
}

func (r *categoryRepository) Update(c *categoryDomain.Category) error {
	ctx := context.Background()

	parsedID, err := uuid.Parse(c.ID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedUserID, err := uuid.Parse(c.UserID)
	if err != nil {
		return err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	// EditedAt をマッピング
	pgEditedAt := pgtype.Timestamptz{Time: c.EditedAt, Valid: true}

	params := dbgen.UpdateCategoryParams{
		Name:     c.Name,
		EditedAt: pgEditedAt,
		ID:       pgID,
		UserID:   pgUserID,
	}

	return r.q.UpdateCategory(ctx, params)
}

func (r *categoryRepository) Delete(categoryID string, userID string) error {
	ctx := context.Background()

	parsedID, err := uuid.Parse(categoryID)
	if err != nil {
		return err
	}
	pgID := pgtype.UUID{Bytes: parsedID, Valid: true}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	params := dbgen.DeleteCategoryParams{
		ID:     pgID,
		UserID: pgUserID,
	}
	return r.q.DeleteCategory(ctx, params)
}
