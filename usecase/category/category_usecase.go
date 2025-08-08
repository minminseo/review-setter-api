package category

import (
	"context"
	"time"

	"github.com/google/uuid"
	categoryDomain "github.com/minminseo/recall-setter/domain/category"
)

type categoryUsecase struct {
	categoryRepo categoryDomain.ICategoryRepository
	// transactionManager transaction.ITransactionManager
}

func NewCategoryUsecase(
	categoryRepo categoryDomain.ICategoryRepository,
	// transactionManager transaction.ITransactionManager,
) ICategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
		// transactionManager: transactionManager,
	}
}

func (cu *categoryUsecase) CreateCategory(ctx context.Context, input CreateCategoryInput) (*CreateCategoryOutput, error) {
	id := uuid.NewString()
	registeredAt := time.Now().UTC()
	editedAt := registeredAt

	newCategory, err := categoryDomain.NewCategory(id, input.UserID, input.Name, registeredAt, editedAt)
	if err != nil {
		return nil, err
	}

	err = cu.categoryRepo.Create(ctx, newCategory)
	if err != nil {
		return nil, err
	}

	resCategory := &CreateCategoryOutput{
		ID:           newCategory.ID(),
		UserID:       newCategory.UserID(),
		Name:         newCategory.Name(),
		RegisteredAt: newCategory.RegisteredAt(),
		EditedAt:     newCategory.EditedAt(),
	}
	return resCategory, nil
}

func (cu *categoryUsecase) GetCategoriesByUserID(ctx context.Context, userID string) ([]*GetCategoryOutput, error) {
	categories, err := cu.categoryRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	outputCategories := make([]*GetCategoryOutput, len(categories))
	for i, c := range categories {
		outputCategories[i] = &GetCategoryOutput{
			ID:           c.ID(),
			UserID:       c.UserID(),
			Name:         c.Name(),
			RegisteredAt: c.RegisteredAt(),
			EditedAt:     c.EditedAt(),
		}
	}
	return outputCategories, nil
}

func (cu *categoryUsecase) UpdateCategory(ctx context.Context, input UpdateCategoryInput) (*UpdateCategoryOutput, error) {
	targetCategory, err := cu.categoryRepo.GetByID(ctx, input.ID, input.UserID)
	if err != nil {
		return nil, err
	}

	EditedAt := time.Now().UTC()

	err = targetCategory.UpdateCategory(input.Name, EditedAt)
	if err != nil {
		return nil, err
	}

	err = cu.categoryRepo.Update(ctx, targetCategory)
	if err != nil {
		return nil, err
	}

	resCategory := &UpdateCategoryOutput{
		ID:       targetCategory.ID(),
		UserID:   targetCategory.UserID(),
		Name:     targetCategory.Name(),
		EditedAt: targetCategory.EditedAt(),
	}
	return resCategory, nil
}

func (cu *categoryUsecase) DeleteCategory(ctx context.Context, categoryID string, userID string) error {
	return cu.categoryRepo.Delete(ctx, categoryID, userID)
}
