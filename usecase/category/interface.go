package category

import "context"

type ICategoryUsecase interface {
	CreateCategory(ctx context.Context, category CreateCategoryInput) (*CreateCategoryOutput, error)
	GetCategoriesByUserID(ctx context.Context, userID string) ([]*GetCategoryOutput, error)
	UpdateCategory(ctx context.Context, category UpdateCategoryInput) (*UpdateCategoryOutput, error)
	DeleteCategory(ctx context.Context, categoryID string, userID string) error
}
