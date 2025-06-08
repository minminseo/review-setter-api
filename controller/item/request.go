package item

type CreateItemRequest struct {
	CategoryID               *string `json:"category_id"`
	BoxID                    *string `json:"box_id"`
	PatternID                *string `json:"pattern_id"`
	Name                     string  `json:"name"`
	Detail                   string  `json:"detail"`
	LearnedDate              string  `json:"learned_date"`
	IsMarkOverdueAsCompleted bool    `json:"is_mark_overdue_as_completed"`
	Today                    string  `json:"today"`
}

// ItemIDはパラメータから取得
type UpdateItemRequest struct {
	CategoryID               *string `json:"category_id"`
	BoxID                    *string `json:"box_id"`
	PatternID                *string `json:"pattern_id"`
	Name                     string  `json:"name"`
	Detail                   string  `json:"detail"`
	LearnedDate              string  `json:"learned_date"`
	IsMarkOverdueAsCompleted bool    `json:"is_mark_overdue_as_completed"`
	Today                    string  `json:"today"`
}

// ReviewDateID、ItemIDはパラメータから取得
type UpdateReviewDatesRequest struct {
	RequestScheduledDate     string                  `json:"request_scheduled_date"`
	IsMarkOverdueAsCompleted bool                    `json:"is_mark_overdue_as_completed"`
	Today                    string                  `json:"today"`
	PatternSteps             []PatternStepForRequest `json:"pattern_steps"`
	LearnedDate              string                  `json:"learned_date"`
	InitialScheduledDate     string                  `json:"initial_scheduled_date"`
	StepNumber               int                     `json:"step_number"`
	CategoryID               *string                 `json:"category_id"`
	BoxID                    *string                 `json:"box_id"`
}

type PatternStepForRequest struct {
	PatternStepID string `json:"pattern_step_id"`
	UserID        string `json:"user_id"`
	PatternID     string `json:"pattern_id"`
	StepNumber    int    `json:"step_number"`
	IntervalDays  int    `json:"interval_days"`
}

type UpdateReviewDateAsCompletedRequest struct {
	StepNumber int `json:"step_number"`
}

type UpdateReviewDateAsInCompletedRequest struct {
	StepNumber int `json:"step_number"`
}
