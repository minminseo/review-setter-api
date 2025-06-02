package box

import "context"

type IBoxUsecase interface {
	CreateBox(ctx context.Context, box CreateBoxInput) (*CreateBoxOutput, error)
	GetBoxesByCategoryID(ctx context.Context, categoryID string, userID string) ([]*GetBoxOutput, error)
	UpdateBox(ctx context.Context, box UpdateBoxInput) (*UpdateBoxOutput, error)
	DeleteBox(ctx context.Context, boxID string, categoryID string, userID string) error
}
