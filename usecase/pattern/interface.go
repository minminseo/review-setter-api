package pattern

import "context"

type IPatternUsecase interface {
	CreatePattern(ctx context.Context, pattern CreatePatternInput) (*CreatePatternOutput, error)
	GetPatternsByUserID(ctx context.Context, userID string) ([]*GetPatternOutput, error)
	UpdatePattern(ctx context.Context, pattern UpdatePatternInput) (*UpdatePatternOutput, error)
	DeletePattern(ctx context.Context, patternID string, userID string) error
}
