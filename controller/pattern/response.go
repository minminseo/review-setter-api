package pattern

import "time"

type PatternStepResponse struct {
	PatternStepID string `json:"pattern_step_id"`
	UserID        string `json:"user_id"`
	PatternID     string `json:"pattern_id"`
	StepNumber    int    `json:"step_number"`
	IntervalDays  int    `json:"interval_days"`
}

type PatternResponse struct {
	ID           string                `json:"id"`
	UserID       string                `json:"user_id"`
	Name         string                `json:"name"`
	TargetWeight string                `json:"target_weight"`
	RegisteredAt time.Time             `json:"registered_at"`
	EditedAt     time.Time             `json:"edited_at"`
	Steps        []PatternStepResponse `json:"steps"`
}
