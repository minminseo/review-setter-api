package item

import "context"

type ItemCountGroupedByBox struct {
	CategoryID string
	BoxID      string
	Count      int
}
type UnclassifiedItemCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type ItemRepository interface {
	Create(ctx context.Context, item *Item) error
	GetDetailByID(ctx context.Context, itemID string, userID string) (*Item, error)
	Update(ctx context.Context, item *Item, userID string) (*Item, error)
	Delete(ctx context.Context, itemID string, userID string) error

	GetAllUnclassifiedByUserID(ctx context.Context, userID string) ([]*Item, error)
	GetAllUnclassifiedByCategoryID(ctx context.Context, categoryID string, userID string) ([]*Item, error)
	GetAllByBoxID(ctx context.Context, boxID string, userID string) ([]*Item, error)

	// ホーム画面では"CountGroupedByBoxByUserID"とCountUnclassifiedByCategoryID"の結果を結合してカテゴリー毎の復習物数を表示
	CountGroupedByBoxByUserID(ctx context.Context, userID string) ([]*ItemCountGroupedByBox, error)
	CountUnclassifiedGroupedByCategoryByUserID(ctx context.Context, userID string) ([]*UnclassifiedItemCountGroupedByCategory, error)

	CountUnclassifiedByUserID(ctx context.Context, serID string) (int, error)
}
