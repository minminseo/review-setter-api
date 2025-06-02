package date

import "context"

type DailyCountGroupedByBox struct {
	CategoryID string
	BoxID      string
	Count      int
}

type UnclassifiedDailyCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type DateRepository interface {
	// itemドメインによる副次的な処理
	Create(ctx context.Context, reviewdates []*Reviewdate) error
	Update(ctx context.Context, reviewdates []*Reviewdate) error

	// ボックス内の復習物毎の復習日一覧を取得
	GetAllByBoxID(ctx context.Context, boxID string, userID string) ([]*Reviewdate, error)
	GetAllUnclassifiedByUserID(ctx context.Context, userID string) ([]*Reviewdate, error)
	GetAllUnclassifiedByCategoryID(ctx context.Context, categoryID string, userID string) ([]*Reviewdate, error)

	// 今日の復習物数の取得系
	// 以下の3つのメソッドで取得した今日の復習物数を組み合わせて、ホーム画面の全体の今日の復習物数を表示
	CountUnclassifiedByUserID(ctx context.Context, userID string) (count int, err error)
	CountGroupedByBoxByUserID(ctx context.Context, userID string) ([]*DailyCountGroupedByBox, error)
	CountUnclassifiedGroupedByCategoryByUserID(ctx context.Context, userID string) ([]*UnclassifiedDailyCountGroupedByCategory, error)
}
