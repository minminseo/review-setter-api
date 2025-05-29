package box

type BoxCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type BoxRepository interface {
	Create(box *Box) error

	// カテゴリー内画面
	// 各ボックス編集画面で表示する情報もこれを使う
	GetAllByCategoryID(categoryID string, userID string) ([]*Box, error)
	GetByID(boxID string, categoryID string, userID string) (*Box, error)

	Update(box *Box) error
	UpdateWithPatternID(box *Box) (int64, error)
	Delete(boxID string, categoryID string, userID string) error

	// TODO: これは復習物リポジトリの責任にする
	// CountGroupedByCategoryByUserID(userID string) ([]*BoxCountGroupedByCategory, error)
}
