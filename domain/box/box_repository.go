package box

import "context"

type BoxCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type IBoxRepository interface {
	Create(ctx context.Context, box *Box) error

	// カテゴリー内画面
	// 各ボックス編集画面で表示する情報もこれを使う
	GetAllByCategoryID(ctx context.Context, categoryID string, userID string) ([]*Box, error)
	GetByID(ctx context.Context, boxID string, categoryID string, userID string) (*Box, error)

	Update(ctx context.Context, box *Box) error
	UpdateWithPatternID(ctx context.Context, box *Box) (int64, error)
	Delete(ctx context.Context, boxID string, categoryID string, userID string) error

	// TODO: これは復習物リポジトリの責任にする
	// CountGroupedByCategoryByUserID(userID string) ([]*BoxCountGroupedByCategory, error)
}
