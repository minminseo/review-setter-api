package box

import (
	"context"
	"time"

	"github.com/google/uuid"
	boxDomain "github.com/minminseo/recall-setter/domain/box"
)

type CreateBoxInput struct {
	UserID     string
	CategoryID string
	PatternID  string
	Name       string
}

type CreateBoxOutput struct {
	ID           string
	UserID       string
	CategoryID   string
	PatternID    string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type GetBoxOutput struct {
	ID           string
	UserID       string
	CategoryID   string
	PatternID    string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type UpdateBoxInput struct {
	ID         string
	UserID     string
	CategoryID string
	PatternID  string
	Name       string
}

type UpdateBoxOutput struct {
	ID         string
	UserID     string
	CategoryID string
	PatternID  string
	Name       string
	EditedAt   time.Time
}

type boxUsecase struct {
	boxRepo boxDomain.IBoxRepository
	// transactionManager transaction.ITransactionManager
}

// NewBoxUsecase はコンストラクタ
func NewBoxUsecase(
	boxRepo boxDomain.IBoxRepository,
	// transactionManager transaction.ITransactionManager,
) IBoxUsecase {
	return &boxUsecase{
		boxRepo: boxRepo,
		// transactionManager: transactionManager,
	}
}

func (bu *boxUsecase) CreateBox(ctx context.Context, input CreateBoxInput) (*CreateBoxOutput, error) {
	id := uuid.NewString()
	registeredAt := time.Now().UTC()
	editedAt := registeredAt

	newBox, err := boxDomain.NewBox(
		id,
		input.UserID,
		input.CategoryID,
		input.PatternID,
		input.Name,
		registeredAt,
		editedAt,
	)
	if err != nil {
		return nil, err
	}

	err = bu.boxRepo.Create(ctx, newBox)
	if err != nil {
		return nil, err
	}

	return &CreateBoxOutput{
		ID:           newBox.ID,
		UserID:       newBox.UserID,
		CategoryID:   newBox.CategoryID,
		PatternID:    newBox.PatternID,
		Name:         newBox.Name,
		RegisteredAt: newBox.RegisteredAt,
		EditedAt:     newBox.EditedAt,
	}, nil
}

func (bu *boxUsecase) GetBoxesByCategoryID(ctx context.Context, categoryID, userID string) ([]*GetBoxOutput, error) {
	boxes, err := bu.boxRepo.GetAllByCategoryID(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}
	outputs := make([]*GetBoxOutput, 0, len(boxes))
	for _, b := range boxes {
		outputs = append(outputs, &GetBoxOutput{
			ID:           b.ID,
			UserID:       b.UserID,
			CategoryID:   b.CategoryID,
			PatternID:    b.PatternID,
			Name:         b.Name,
			RegisteredAt: b.RegisteredAt,
			EditedAt:     b.EditedAt,
		})
	}
	return outputs, nil
}

func (bu *boxUsecase) UpdateBox(ctx context.Context, input UpdateBoxInput) (*UpdateBoxOutput, error) {
	targetBox, err := bu.boxRepo.GetByID(ctx, input.ID, input.CategoryID, input.UserID)
	if err != nil {
		return nil, err
	}

	editedAt := time.Now().UTC()

	var isSamePattern bool
	isSamePattern, err = targetBox.Set(input.PatternID, input.Name, editedAt)
	if err != nil {
		return nil, err
	}

	if isSamePattern {
		err = bu.boxRepo.Update(ctx, targetBox)
		if err != nil {
			return nil, err
		}
	} else {
		affected, err := bu.boxRepo.UpdateWithPatternID(ctx, targetBox)
		if err != nil {
			return nil, err
		}
		if affected == 0 {
			return nil, boxDomain.ErrPatternConflict
		}
	}

	resBox := &UpdateBoxOutput{
		ID:         targetBox.ID,
		UserID:     targetBox.UserID,
		CategoryID: targetBox.CategoryID,
		PatternID:  targetBox.PatternID,
		Name:       targetBox.Name,
		EditedAt:   targetBox.EditedAt,
	}
	return resBox, nil
}

func (uc *boxUsecase) DeleteBox(ctx context.Context, boxID string, categoryID string, userID string) error {
	return uc.boxRepo.Delete(ctx, boxID, categoryID, userID)
}
