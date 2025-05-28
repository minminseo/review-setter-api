package category

import (
	"github.com/google/uuid"
	categoryDomain "github.com/minminseo/recall-setter/domain/category"
)

type CreateCategoryInput struct {
	UserID string
	Name   string
}

type CreateCategoryOutput struct {
	ID     string
	UserID string
	Name   string
}

type GetCategoryOutput struct {
	ID     string
	UserID string
	Name   string
}

type UpdateCategoryInput struct {
	ID     string
	UserID string
	Name   string
}

type UpdateCategoryOutput struct {
	ID     string
	UserID string
	Name   string
}

type categoryUsecase struct {
	categoryRepo categoryDomain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo categoryDomain.CategoryRepository) ICategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (cu *categoryUsecase) CreateCategory(input CreateCategoryInput) (*CreateCategoryOutput, error) {
	id := uuid.NewString()

	newCategory, err := categoryDomain.NewCategory(id, input.UserID, input.Name)
	if err != nil {
		return nil, err
	}

	err = cu.categoryRepo.Create(newCategory)
	if err != nil {
		return nil, err
	}

	resCategory := &CreateCategoryOutput{
		ID:     newCategory.ID,
		UserID: newCategory.UserID,
		Name:   newCategory.Name,
	}
	return resCategory, nil
}

func (cu *categoryUsecase) GetCategoriesByUserID(userID string) ([]*GetCategoryOutput, error) {
	categories, err := cu.categoryRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var outputCategories []*GetCategoryOutput
	for _, cat := range categories {
		outputCategories = append(outputCategories, &GetCategoryOutput{
			ID:     cat.ID,
			UserID: cat.UserID,
			Name:   cat.Name,
		})
	}
	return outputCategories, nil
}

func (cu *categoryUsecase) UpdateCategory(input UpdateCategoryInput) (*UpdateCategoryOutput, error) {
	targetCategory, err := cu.categoryRepo.GetByID(input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	err = targetCategory.Set(input.Name)
	if err != nil {
		return nil, err
	}

	err = cu.categoryRepo.Update(targetCategory)
	if err != nil {
		return nil, err
	}

	resCategory := &UpdateCategoryOutput{
		ID:     targetCategory.ID,
		UserID: targetCategory.UserID,
		Name:   targetCategory.Name,
	}
	return resCategory, nil
}

func (cu *categoryUsecase) DeleteCategory(categoryID string, userID string) error {
	return cu.categoryRepo.Delete(categoryID, userID)
}
