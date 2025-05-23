package item

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
	Create(item *Item) error
	GetDetailByID(itemID string, userID string) (*Item, error)
	Update(item *Item, userID string) (*Item, error)
	Delete(itemID string, userID string) error

	GetAllUnclassifiedByUserID(userID string) ([]*Item, error)
	GetAllUnclassifiedByCategoryID(categoryID string, userID string) ([]*Item, error)
	GetAllByBoxID(boxID string, userID string) ([]*Item, error)

	// ホーム画面では"CountGroupedByBoxByUserID"とCountUnclassifiedByCategoryID"の結果を結合してカテゴリー毎の復習物数を表示
	CountGroupedByBoxByUserID(userID string) ([]*ItemCountGroupedByBox, error)
	CountUnclassifiedGroupedByCategoryByUserID(userID string) ([]*UnclassifiedItemCountGroupedByCategory, error)

	CountUnclassifiedByUserID(userID string) (int, error)
}
