package box

type BoxCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type BoxRepository interface {
	Create(box *box) (*box, error)
	Update(box *box, userID string) (*box, error)
	Delete(boxID string, userID string) error

	// カテゴリー内画面
	// 各ボックス編集画面で表示する情報もこれを使う
	GetAllByCategoryID(categoryID string, userID string) ([]*box, error)

	CountGroupedByCategoryByUserID(userID string) ([]*BoxCountGroupedByCategory, error)
}
