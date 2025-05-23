package category

type CategoryRepository interface {
	Create(category *Category) (*Category, error)
	GetAllByUserID(userID string) ([]*Category, error)
	Update(category *Category, userID string) (*Category, error)
	Delete(categoryID string, userID string) error
}
