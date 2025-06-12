package item

import "context"

type IItemUsecase interface {
	CreateItem(ctx context.Context, item CreateItemInput) (*CreateItemOutput, error)
	UpdateItem(ctx context.Context, item UpdateItemInput) (*UpdateItemOutput, error)
	UpdateReviewDates(ctx context.Context, input UpdateBackReviewDateInput) (*UpdateBackReviewDateOutput, error)
	UpdateItemAsFinishedForce(ctx context.Context, input UpdateItemAsFinishedForceInput) (*UpdateItemAsFinishedForceOutput, error)
	UpdateReviewDateAsCompleted(ctx context.Context, input UpdateReviewDateAsCompletedInput) (*UpdateReviewDateAsCompletedOutput, error)
	UpdateReviewDateAsInCompleted(ctx context.Context, input UpdateReviewDateAsInCompletedInput) (*UpdateReviewDateAsInCompletedOutput, error)
	UpdateItemAsUnFinishedForce(ctx context.Context, input UpdateItemAsUnFinishedForceInput) (*UpdateItemAsUnFinishedForceOutput, error)
	DeleteItem(ctx context.Context, itemID string, userID string) error

	/* ボックス内の復習物一覧表示のための取得メソッド*/
	GetAllUnFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*GetItemOutput, error)
	GetAllUnFinishedUnclassifiedItemsByUserID(ctx context.Context, userID string) ([]*GetItemOutput, error)
	GetAllUnFinishedUnclassifiedItemsByCategoryID(ctx context.Context, userID string, categoryID string) ([]*GetItemOutput, error)

	// アプリ内に存在するデータたちの概要を表示するための取得メソッド
	// 復習物数系
	CountItemsGroupedByBoxByUserID(ctx context.Context, userID string) ([]*ItemCountGroupedByBoxOutput, error)
	CountUnclassifiedItemsGroupedByCategoryByUserID(ctx context.Context, userID string) ([]*UnclassifiedItemCountGroupedByCategoryOutput, error)
	CountUnclassifiedItemsByUserID(ctx context.Context, userID string) (int, error)

	// 今日の復習物（復習日）数系
	CountDailyDatesGroupedByBoxByUserID(ctx context.Context, userID string, today string) ([]*DailyCountGroupedByBoxOutput, error)
	CountDailyDatesUnclassifiedGroupedByCategoryByUserID(ctx context.Context, userID string, today string) ([]*UnclassifiedDailyDatesCountGroupedByCategoryOutput, error)
	CountDailyDatesUnclassifiedByUserID(ctx context.Context, userID string, today string) (int, error)

	// 今日の全復習日数を取得する
	CountAllDailyReviewDates(ctx context.Context, userID string, today string) (int, error)

	// 今日の復習日一覧を取得する
	GetAllDailyReviewDates(ctx context.Context, userID string, today string) (*GetDailyReviewDatesOutput, error)

	// 完了済み復習物を取得する系
	GetFinishedItemsByBoxID(ctx context.Context, boxID string, userID string) ([]*GetItemOutput, error)
	GetUnclassfiedFinishedItemsByCategoryID(ctx context.Context, userID string, categoryID string) ([]*GetItemOutput, error)
	GetUnclassfiedFinishedItemsByUserID(ctx context.Context, userID string) ([]*GetItemOutput, error)
}
