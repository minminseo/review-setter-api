package category

type CategoryRepository interface {
	Create(category *Category) error
	GetAllByUserID(userID string) ([]*Category, error)
	GetByID(categoryID string, userID string) (*Category, error)
	Update(category *Category) error
	Delete(categoryID string, userID string) error
}
