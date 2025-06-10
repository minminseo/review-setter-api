package box

import "context"

type BoxCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type BoxName struct {
	BoxID     string
	Name      string
	PatternID string
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

	// item_usecaseで使う。ボックスの名前とパターンIDを一覧取得する
	GetBoxNamesByBoxIDs(ctx context.Context, boxIDs []string) ([]*BoxName, error)
}
