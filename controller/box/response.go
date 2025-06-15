package box

import "time"

type BoxResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	CategoryID   string    `json:"category_id"`
	PatternID    string    `json:"pattern_id"`
	Name         string    `json:"name"`
	RegisteredAt time.Time `json:"registered_at"`
	EditedAt     time.Time `json:"edited_at"`
}

type UpdateBoxResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CategoryID string    `json:"category_id"`
	PatternID  string    `json:"pattern_id"`
	Name       string    `json:"name"`
	EditedAt   time.Time `json:"edited_at"`
}
