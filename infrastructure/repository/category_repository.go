package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	categoryDomain "github.com/minminseo/recall-setter/domain/category"
	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
)

type categoryRepository struct{}

func NewCategoryRepository() categoryDomain.ICategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) Create(ctx context.Context, category *categoryDomain.Category) error {
	q := db.GetQuery(ctx)

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

	pgRegisteredAt := pgtype.Timestamptz{Time: category.RegisteredAt, Valid: true}
	pgEditedAt := pgtype.Timestamptz{Time: category.EditedAt, Valid: true}

	params := dbgen.CreateCategoryParams{
		ID:           pgID,
		UserID:       pgUserID,
		Name:         category.Name,
		RegisteredAt: pgRegisteredAt,
		EditedAt:     pgEditedAt,
	}

	return q.CreateCategory(ctx, params)
}

func (r *categoryRepository) GetAllByUserID(ctx context.Context, userID string) ([]*categoryDomain.Category, error) {
	q := db.GetQuery(ctx)

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pgUserID := pgtype.UUID{Bytes: parsedUserID, Valid: true}

	rows, err := q.GetAllCategoriesByUserID(ctx, pgUserID)
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

func (r *categoryRepository) GetByID(ctx context.Context, categoryID string, userID string) (*categoryDomain.Category, error) {
	q := db.GetQuery(ctx)

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
	row, err := q.GetCategoryByID(ctx, params)
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

func (r *categoryRepository) Update(ctx context.Context, category *categoryDomain.Category) error {
	q := db.GetQuery(ctx)

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

	pgEditedAt := pgtype.Timestamptz{Time: category.EditedAt, Valid: true}

	params := dbgen.UpdateCategoryParams{
		Name:     category.Name,
		EditedAt: pgEditedAt,
		ID:       pgID,
		UserID:   pgUserID,
	}
	return q.UpdateCategory(ctx, params)
}

func (r *categoryRepository) Delete(ctx context.Context, categoryID string, userID string) error {
	q := db.GetQuery(ctx)

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
	return q.DeleteCategory(ctx, params)
}
