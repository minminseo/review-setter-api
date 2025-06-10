package category

import "context"

type CategoryName struct {
	ID   string
	Name string
}

type ICategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetAllByUserID(ctx context.Context, userID string) ([]*Category, error)
	GetByID(ctx context.Context, categoryID string, userID string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, categoryID string, userID string) error

	// item_usecaseで使う。カテゴリーの名前を一覧取得する
	GetCategoryNamesByCategoryIDs(ctx context.Context, categoryIDs []string) ([]*CategoryName, error)
}
