package item

import (
	"context"
	"time"
)

type ItemCountGroupedByBox struct {
	CategoryID string
	BoxID      string
	Count      int
}
type UnclassifiedItemCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type DailyCountGroupedByBox struct {
	CategoryID string
	BoxID      string
	Count      int
}

type UnclassifiedDailyDatesCountGroupedByCategory struct {
	CategoryID string
	Count      int
}

type DailyReviewDate struct {
	ReviewdateID      string
	CategoryID        *string
	BoxID             *string
	StepNumber        int
	PrevScheduledDate *time.Time
	ScheduledDate     time.Time
	NextScheduledDate *time.Time
	IsCompleted       bool
	ItemID            string
	Name              string
	Detail            string
	LearnedDate       time.Time
	RegisteredAt      time.Time
	EditedAt          time.Time
}

type IItemRepository interface {
	CreateItem(ctx context.Context, item *Item) error
	CreateReviewdates(ctx context.Context, reviewdates []*Reviewdate) (int64, error)

	// 復習パターンor復習日が変更対象なのかどうか判定するためのメソッド
	GetItemByID(ctx context.Context, itemID string, userID string) (*Item, error)

	// 復習物が持つ復習日を取得して、完了済みの復習日がないか判別するためのメソッド
	HasCompletedReviewDateByItemID(ctx context.Context, itemID string, userID string) (bool, error)

	GetReviewDateIDsByItemID(ctx context.Context, itemID string, userID string) ([]string, error)

	// 復習物の更新
	UpdateItem(ctx context.Context, item *Item) error

	// 復習日の更新。
	// //この更新には、「復習パターン変更による更新」、「ボックス移動（復習パターン不一致時）による更新」は含まれない（代わりにCreateReviewdatesを使う）
	UpdateReviewDates(ctx context.Context, reviewdates []*Reviewdate, userID string) error
	UpdateReviewDatesBack(ctx context.Context, reviewdates []*Reviewdate, userID string) error

	// 復習物の途中完了（手動）の場合 or 復習日巻き戻で全完了した場合 or 通常の復習日完了操作による自動完了
	UpdateItemAsFinished(ctx context.Context, itemID string, userID string, editedAt time.Time) error

	UpdateItemAsUnFinished(ctx context.Context, itemID string, userID string, editedAt time.Time) error

	// 復習日を完了済みに更新
	UpdateReviewDateAsCompleted(ctx context.Context, reviewdateID string, userID string) error

	UpdateReviewDateAsInCompleted(ctx context.Context, reviewdateID string, userID string) error

	// 復習日巻き戻し操作時の最新復習スケジュールを取得するため・復習日完了操作対象の復習日が最後の復習日かどうか判別するため
	GetReviewDatesByItemID(ctx context.Context, itemID string, userID string) ([]*Reviewdate, error)

	// 復習物の削除
	DeleteItem(ctx context.Context, itemID string, userID string) error

	DeleteReviewDates(ctx context.Context, itemID string, userID string) error

	/*-------------*/
	// ここからしたは取得系

	// ボックスAの復習物（未完了）と復習日を一覧取得
	GetAllUnFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*Item, error)

	// 完了済みの復習物に紐づいた復習日もまとめて取得。完了済み復習物が持つ復習日を除外する必要がある。
	// TODO: rview_datesテーブルにも親の復習物が完了かどうかのフラグを持たせるべきか検討（冗長化）。
	GetAllReviewDatesByBoxID(ctx context.Context, boxID string, userID string) ([]*Reviewdate, error)

	// ホーム画面の未分類復習物ボックスの復習物（未完了）と復習日を一覧取得
	GetAllUnFinishedUnclassifiedItemsByUserID(ctx context.Context, userID string) ([]*Item, error)
	GetAllUnclassifiedReviewDatesByUserID(ctx context.Context, userID string) ([]*Reviewdate, error)

	// カテゴリーAの未分類復習物ボックスの復習物（未完了）と復習日を一覧取得
	GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx context.Context, categoryID string, userID string) ([]*Item, error)
	GetAllUnclassifiedReviewDatesByCategoryID(ctx context.Context, categoryID string, userID string) ([]*Reviewdate, error)

	/*--------------------------------------*/

	//ここから下は概要表示用の取得メソッド
	/*--------------------------------------*/

	// 復習物数系
	// ホーム画面では"CountItemsGroupedByBoxByUserID"と"CountUnclassifiedItemsGroupedByCategoryByUserID"の結果を結合してカテゴリー毎の復習物数を表示
	// カテゴリー毎の全復習物ボックス毎の復習物数を取得
	CountItemsGroupedByBoxByUserID(ctx context.Context, userID string) ([]*ItemCountGroupedByBox, error)

	// カテゴリー毎の未分類復習物ボックスの復習物数を取得
	CountUnclassifiedItemsGroupedByCategoryByUserID(ctx context.Context, userID string) ([]*UnclassifiedItemCountGroupedByCategory, error)

	// ホーム画面の未分類復習物ボックスの復習物数を取得
	CountUnclassifiedItemsByUserID(ctx context.Context, userID string) (int, error)

	/*--------------------------------------*/

	// 今日の復習物（復習日）系
	// 以下の3つのメソッドで取得した今日の復習物数を組み合わせて、ホーム画面の全体の今日の復習物数を表示
	// カテゴリー毎の全復習物ボックス毎の今日の復習物数（復習日）を取得
	CountDailyDatesGroupedByBoxByUserID(ctx context.Context, userID string, targetDate time.Time) ([]*DailyCountGroupedByBox, error)

	// カテゴリー毎の未分類復習物ボックスの今日の復習物数（復習日）を取得
	CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx context.Context, userID string, targetDate time.Time) ([]*UnclassifiedDailyDatesCountGroupedByCategory, error)

	// ホーム画面の未分類復習物ボックスの今日の復習物数（復習日）を取得
	CountDailyDatesUnclassifiedByUserID(ctx context.Context, userID string, targetDate time.Time) (int, error)

	// EditedAtの取得専用
	GetEditedAtByItemID(ctx context.Context, itemID string, userID string) (time.Time, error)

	// 今日の全復習日数を取得する
	CountAllDailyReviewDates(ctx context.Context, userID string, parsedToday time.Time) (int, error)

	GetAllDailyReviewDates(ctx context.Context, userID string, parsedToday time.Time) ([]*DailyReviewDate, error)

	// 完了済み復習物系を取得する系
	GetFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*Item, error)
	GetUnclassfiedFinishedItemsByCategoryID(ctx context.Context, categoryID string, userID string) ([]*Item, error)
	GetUnclassfiedFinishedItemsByUserID(ctx context.Context, userID string) ([]*Item, error)

	/*--------------------*/
	// patternパッケージで使うメソッド
	IsPatternRelatedToItemByPatternID(ctx context.Context, patternID string, userID string) (bool, error)
}
