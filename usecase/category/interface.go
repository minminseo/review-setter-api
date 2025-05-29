package category

type ICategoryUsecase interface {
	CreateCategory(category CreateCategoryInput) (*CreateCategoryOutput, error)
	GetCategoriesByUserID(userID string) ([]*GetCategoryOutput, error)
	UpdateCategory(category UpdateCategoryInput) (*UpdateCategoryOutput, error)
	DeleteCategory(categoryID string, userID string) error
}
