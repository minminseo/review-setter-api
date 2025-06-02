package category

import "context"

type ICategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetAllByUserID(ctx context.Context, userID string) ([]*Category, error)
	GetByID(ctx context.Context, categoryID string, userID string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, categoryID string, userID string) error
}
