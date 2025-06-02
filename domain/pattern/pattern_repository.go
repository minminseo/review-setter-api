package pattern

import (
	"context"
)

type IPatternRepository interface {
	CreatePattern(ctx context.Context, pattern *Pattern) error
	CreatePatternSteps(ctx context.Context, steps []*PatternStep) (int64, error)

	// 復習パターン一覧取得機能用
	GetAllPatternsByUserID(ctx context.Context, userID string) ([]*Pattern, error)
	GetAllPatternStepsByUserID(ctx context.Context, userID string) ([]*PatternStep, error)

	UpdatePattern(ctx context.Context, pattern *Pattern) error
	DeletePattern(ctx context.Context, patternID string, userID string) error
	DeletePatternSteps(ctx context.Context, patternID string, userID string) error

	// ボックス一覧取得→ボックス毎にループ処理（Patternを取得→PatternStepたちを取得）
	FindPatternByPatternID(ctx context.Context, patternID string, userID string) (*Pattern, error)
	GetAllPatternStepsByPatternID(ctx context.Context, patternID string, userID string) ([]*PatternStep, error)
}
