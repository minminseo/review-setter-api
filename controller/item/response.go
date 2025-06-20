package item

import "time"

type ReviewDateResponse struct {
	ReviewDateID         string  `json:"review_date_id"`
	UserID               string  `json:"user_id"`
	CategoryID           *string `json:"category_id"`
	BoxID                *string `json:"box_id"`
	ItemID               string  `json:"item_id"`
	StepNumber           int     `json:"step_number"`
	InitialScheduledDate string  `json:"initial_scheduled_date"`
	ScheduledDate        string  `json:"scheduled_date"`
	IsCompleted          bool    `json:"is_completed"`
}

type ItemResponse struct {
	ItemID       string               `json:"item_id"`
	UserID       string               `json:"user_id"`
	CategoryID   *string              `json:"category_id"`
	BoxID        *string              `json:"box_id"`
	PatternID    *string              `json:"pattern_id"`
	Name         string               `json:"name"`
	Detail       string               `json:"detail"`
	LearnedDate  string               `json:"learned_date"`
	IsFinished   bool                 `json:"is_finished"`
	RegisteredAt time.Time            `json:"registered_at"`
	EditedAt     time.Time            `json:"edited_at"`
	ReviewDates  []ReviewDateResponse `json:"review_dates"`
}

type UpdateItemAsFinishedForceResponse struct {
	ItemID     string    `json:"item_id"`
	UserID     string    `json:"user_id"`
	IsFinished bool      `json:"is_finished"`
	EditedAt   time.Time `json:"edited_at"`
}

type UpdateReviewDateAsCompletedResponse struct {
	ReviewDateID string    `json:"review_date_id"`
	UserID       string    `json:"user_id"`
	IsCompleted  bool      `json:"is_completed"`
	IsFinished   bool      `json:"is_finished"`
	EditedAt     time.Time `json:"edited_at"`
}

type UpdateReviewDateAsInCompletedResponse struct {
	ReviewDateID string    `json:"review_date_id"`
	UserID       string    `json:"user_id"`
	IsCompleted  bool      `json:"is_completed"`
	IsFinished   bool      `json:"is_finished"`
	EditedAt     time.Time `json:"edited_at"`
}

type UpdateBackReviewDateResponse struct {
	ItemID      string               `json:"item_id"`
	UserID      string               `json:"user_id"`
	IsFinished  bool                 `json:"is_finished"`
	EditedAt    time.Time            `json:"edited_at"`
	ReviewDates []ReviewDateResponse `json:"review_dates"`
}

type UpdateItemAsUnFinishedForceResponse struct {
	ItemID      string               `json:"item_id"`
	UserID      string               `json:"user_id"`
	IsFinished  bool                 `json:"is_finished"`
	EditedAt    time.Time            `json:"edited_at"`
	ReviewDates []ReviewDateResponse `json:"review_dates"`
}

type CountResponse struct {
	Count int `json:"count"`
}

type ItemCountGroupedByBoxResponse struct {
	CategoryID string `json:"category_id"`
	BoxID      string `json:"box_id"`
	Count      int    `json:"count"`
}

type UnclassifiedItemCountGroupedByCategoryResponse struct {
	CategoryID string `json:"category_id"`
	Count      int    `json:"count"`
}

type DailyCountGroupedByBoxResponse struct {
	CategoryID string `json:"category_id"`
	BoxID      string `json:"box_id"`
	Count      int    `json:"count"`
}

type UnclassifiedDailyDatesCountGroupedByCategoryResponse struct {
	CategoryID string `json:"category_id"`
	Count      int    `json:"count"`
}

type DailyReviewDatesByBoxResponse struct {
	ReviewDateID      string    `json:"review_date_id"`
	CategoryID        string    `json:"category_id"`
	BoxID             string    `json:"box_id"`
	StepNumber        int       `json:"step_number"`
	PrevScheduledDate *string   `json:"prev_scheduled_date"`
	ScheduledDate     string    `json:"scheduled_date"`
	NextScheduledDate *string   `json:"next_scheduled_date"`
	IsCompleted       bool      `json:"is_completed"`
	ItemName          string    `json:"item_name"`
	Detail            string    `json:"detail"`
	LearnedDate       string    `json:"learned_date"`
	RegisteredAt      time.Time `json:"registered_at"`
	EditedAt          time.Time `json:"edited_at"`
}

type DailyReviewDatesGroupedByBoxResponse struct {
	BoxID        string                          `json:"box_id"`
	CategoryID   string                          `json:"category_id"`
	BoxName      string                          `json:"box_name"`
	ReviewDates  []DailyReviewDatesByBoxResponse `json:"review_dates"`
	TargetWeight string                          `json:"target_weight"`
}

type UnclassifiedDailyReviewDatesGroupedByCategoryResponse struct {
	ReviewDateID      string    `json:"review_date_id"`
	CategoryID        string    `json:"category_id"`
	StepNumber        int       `json:"step_number"`
	PrevScheduledDate *string   `json:"prev_scheduled_date"`
	ScheduledDate     string    `json:"scheduled_date"`
	NextScheduledDate *string   `json:"next_scheduled_date"`
	IsCompleted       bool      `json:"is_completed"`
	ItemName          string    `json:"item_name"`
	Detail            string    `json:"detail"`
	LearnedDate       string    `json:"learned_date"`
	RegisteredAt      time.Time `json:"registered_at"`
	EditedAt          time.Time `json:"edited_at"`
}

type DailyReviewDatesGroupedByCategoryResponse struct {
	CategoryID                             string                                                  `json:"category_id"`
	CategoryName                           string                                                  `json:"category_name"`
	Boxes                                  []DailyReviewDatesGroupedByBoxResponse                  `json:"boxes"`
	UnclassifiedDailyReviewDatesByCategory []UnclassifiedDailyReviewDatesGroupedByCategoryResponse `json:"unclassified_daily_review_dates_by_category"`
}

type UnclassifiedDailyReviewDatesGroupedByUserResponse struct {
	ReviewDateID      string    `json:"review_date_id"`
	StepNumber        int       `json:"step_number"`
	PrevScheduledDate *string   `json:"prev_scheduled_date"`
	ScheduledDate     string    `json:"scheduled_date"`
	NextScheduledDate *string   `json:"next_scheduled_date"`
	IsCompleted       bool      `json:"is_completed"`
	ItemName          string    `json:"item_name"`
	Detail            string    `json:"detail"`
	LearnedDate       string    `json:"learned_date"`
	RegisteredAt      time.Time `json:"registered_at"`
	EditedAt          time.Time `json:"edited_at"`
}

type GetDailyReviewDatesResponse struct {
	Categories                    []DailyReviewDatesGroupedByCategoryResponse         `json:"categories"`
	DailyReviewDatesGroupedByUser []UnclassifiedDailyReviewDatesGroupedByUserResponse `json:"daily_review_dates_grouped_by_user"`
}
