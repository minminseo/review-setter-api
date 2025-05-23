package date

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
	Create(reviewdates []*Reviewdate) error
	Update(reviewdates []*Reviewdate) error

	// ボックス内の復習物毎の復習日一覧を取得
	GetAllByBoxID(boxID string, userID string) ([]*Reviewdate, error)
	GetAllUnclassifiedByUserID(userID string) ([]*Reviewdate, error)
	GetAllUnclassifiedByCategoryID(categoryID string, userID string) ([]*Reviewdate, error)

	// 今日の復習物数の取得系
	// 以下の3つのメソッドで取得した今日の復習物数を組み合わせて、ホーム画面の全体の今日の復習物数を表示
	CountUnclassifiedByUserID(userID string) (count int, err error)
	CountGroupedByBoxByUserID(userID string) ([]*DailyCountGroupedByBox, error)
	CountUnclassifiedGroupedByCategoryByUserID(userID string) ([]*UnclassifiedDailyCountGroupedByCategory, error)
}
