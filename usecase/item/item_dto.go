package item

import "time"

type CreateItemInput struct {
	UserID                   string
	CategoryID               *string
	BoxID                    *string
	PatternID                *string
	Name                     string
	Detail                   string
	LearnedDate              string
	IsMarkOverdueAsCompleted bool
	Today                    string
}

type CreateReviewdateOutput struct {
	DateID               string
	UserID               string
	ItemID               string
	StepNumber           int
	InitialScheduledDate string
	ScheduledDate        string
	IsCompleted          bool
}

type CreateItemOutput struct {
	ItemID       string
	UserID       string
	CategoryID   *string
	BoxID        *string
	PatternID    *string
	Name         string
	Detail       string
	LearnedDate  string
	IsCompleted  bool
	RegisteredAt time.Time
	EditedAt     time.Time
	Reviewdates  []CreateReviewdateOutput
}

// 分類済みボックスの復習物の更新リクエスト用のDTO

type UpdateItemInput struct {
	ItemID                   string
	UserID                   string
	CategoryID               *string
	BoxID                    *string
	PatternID                *string
	Name                     string
	Detail                   string
	LearnedDate              string
	IsMarkOverdueAsCompleted bool
	Today                    string
}

// 復習物更新用の分類、未分類共通DTO
type UpdateReviewDateOutput struct {
	ReviewDateID         string
	UserID               string
	CategoryID           *string
	BoxID                *string
	ItemID               string
	StepNumber           int
	InitialScheduledDate string
	ScheduledDate        string
	IsCompleted          bool
}

type UpdateItemOutput struct {
	ItemID      string
	UserID      string
	CategoryID  *string
	BoxID       *string
	PatternID   *string
	Name        string
	Detail      string
	LearnedDate string
	IsFinished  bool
	EditedAt    time.Time
	ReviewDates []UpdateReviewDateOutput
}

// 復習物の途中完了（手動）リクエスト用のDTO
type UpdateItemAsFinishedForceInput struct {
	ItemID string
	UserID string
}

type UpdateItemAsFinishedForceOutput struct {
	ItemID     string
	UserID     string
	IsFinished bool
	EditedAt   time.Time
}

// 途中完了した復習物の再開リクエスト用のDTO
type UpdateItemAsUnFinishedForceInput struct {
	ItemID      string
	UserID      string
	CategoryID  *string // nilなら未分類
	BoxID       *string // nilなら未分類
	PatternID   string  // 完了済みボックスの時点で必ずある
	LearnedDate string  // 計算用
	Today       string
}

type UpdateItemAsUnFinishedForceOutput struct {
	ItemID      string
	UserID      string
	IsFinished  bool
	EditedAt    time.Time
	ReviewDates []UpdateReviewDateOutput // 復習物更新のDTO共有
}

type UpdateReviewDateAsCompletedInput struct {
	ReviewDateID string
	UserID       string
	ItemID       string
	StepNumber   int
}

// 全ての復習日が完了したかどうかも返す（IsFinished）
type UpdateReviewDateAsCompletedOutput struct {
	ReviewDateID string
	UserID       string
	IsCompleted  bool
	IsFinished   bool
	EditedAt     time.Time
}

type UpdateReviewDateAsInCompletedInput struct {
	ReviewDateID string
	UserID       string
	ItemID       string
	StepNumber   int
}

type UpdateReviewDateAsInCompletedOutput struct {
	ReviewDateID string
	UserID       string
	IsCompleted  bool
	IsFinished   bool
	EditedAt     time.Time
}

type PatternStepInReviewDate struct {
	PatternStepID string
	UserID        string
	PatternID     string
	StepNumber    int
	IntervalDays  int
}

// 復習日の変更（今日よりも前に巻き戻す）は未分類かどうかは区別しない
type UpdateBackReviewDateInput struct {
	ReviewDateID             string
	UserID                   string
	CategoryID               *string // nilなら未分類
	BoxID                    *string // nilなら未分類
	ItemID                   string
	StepNumber               int
	InitialScheduledDate     string
	RequestScheduledDate     string
	IsMarkOverdueAsCompleted bool
	Today                    string

	// 計算用
	LearnedDate string // 復習物の学習日

	// 復習日再生成用（DBアクセス減らす目的）
	PatternStepsInReviewDate []PatternStepInReviewDate // 復習物更新のDTO共有（）
}

type UpdateBackReviewDateOutput struct {
	ItemID      string
	UserID      string
	IsFinished  bool
	EditedAt    time.Time
	ReviewDates []UpdateReviewDateOutput // 復習物更新のDTO共有
}

/* --- */
/* ボックス内の復習物一覧表示のための取得メソッド*/

type GetReviewDateOutput struct {
	ReviewDateID         string
	UserID               string
	CategoryID           *string // nilなら未分類
	BoxID                *string // nilなら未分類
	ItemID               string
	StepNumber           int
	InitialScheduledDate string
	ScheduledDate        string
	IsCompleted          bool
}

type GetItemOutput struct {
	ItemID       string
	UserID       string
	CategoryID   *string
	BoxID        *string
	PatternID    *string
	Name         string
	Detail       string
	LearnedDate  string
	IsFinished   bool
	RegisteredAt time.Time
	EditedAt     time.Time
	ReviewDates  []GetReviewDateOutput // ポインタにしたらどうなる？
}

// アプリ内に存在するデータたちの概要を表示するための取得メソッド
type ItemCountGroupedByBoxOutput struct {
	CategoryID string
	BoxID      string
	Count      int
}

type UnclassifiedItemCountGroupedByCategoryOutput struct {
	CategoryID string
	Count      int
}

type DailyCountGroupedByBoxOutput struct {
	CategoryID string
	BoxID      string
	Count      int
}

type UnclassifiedDailyDatesCountGroupedByCategoryOutput struct {
	CategoryID string
	Count      int
}

/*
変更したい日付の変更前が今日じゃないならエラー→これTZ考慮必要じゃん
変更後の日付がinitial_scheduled_dateより前ならエラー


*/

type DailyReviewDatesByBoxOutput struct {
	ReviewDateID      string
	CategoryID        string // 必ず存在する
	BoxID             string // 必ず存在する
	StepNumber        int
	PrevScheduledDate *string
	ScheduledDate     string
	NextScheduledDate *string
	IsCompleted       bool

	// 復習物の情報
	ItemName     string
	Detail       string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type DailyReviewDatesGroupedByBoxOutput struct {
	BoxID        string
	CategoryID   string
	BoxName      string
	ReviewDates  []DailyReviewDatesByBoxOutput
	TargetWeight string
}

type UnclassifiedDailyReviewDatesGroupedByCategoryOutput struct {
	ReviewDateID      string
	CategoryID        string // 必ず存在する
	StepNumber        int
	PrevScheduledDate *string
	ScheduledDate     string
	NextScheduledDate *string
	IsCompleted       bool

	// 復習物の情報
	ItemName     string
	Detail       string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type DailyReviewDatesGroupedByCategoryOutput struct {
	CategoryID                             string
	CategoryName                           string
	Boxes                                  []DailyReviewDatesGroupedByBoxOutput
	UnclassifiedDailyReviewDatesByCategory []UnclassifiedDailyReviewDatesGroupedByCategoryOutput
}

type UnclassifiedDailyReviewDatesGroupedByUserOutput struct {
	ReviewDateID      string
	StepNumber        int
	PrevScheduledDate *string
	ScheduledDate     string
	NextScheduledDate *string
	IsCompleted       bool

	// 復習物の情報
	ItemName     string
	Detail       string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type GetDailyReviewDatesOutput struct {
	Categories                    []DailyReviewDatesGroupedByCategoryOutput
	DailyReviewDatesGroupedByUser []UnclassifiedDailyReviewDatesGroupedByUserOutput
}
