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
	RequestScheduledDate     string  `json:"request_scheduled_date"`
	IsMarkOverdueAsCompleted bool    `json:"is_mark_overdue_as_completed"`
	Today                    string  `json:"today"`
	PatternID                string  `json:"pattern_id"`
	LearnedDate              string  `json:"learned_date"`
	InitialScheduledDate     string  `json:"initial_scheduled_date"`
	StepNumber               int     `json:"step_number"`
	CategoryID               *string `json:"category_id"`
	BoxID                    *string `json:"box_id"`
}

type UpdateItemAsUnFinishedForceRequest struct {
	CategoryID  *string `json:"category_id"`
	BoxID       *string `json:"box_id"`
	PatternID   string  `json:"pattern_id"`
	LearnedDate string  `json:"learned_date"`
	Today       string  `json:"today"`
}

type UpdateReviewDateAsCompletedRequest struct {
	StepNumber int `json:"step_number"`
}

type UpdateReviewDateAsInCompletedRequest struct {
	StepNumber int `json:"step_number"`
}
