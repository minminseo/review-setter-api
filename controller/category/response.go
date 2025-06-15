package category

import "time"

type CategoryResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	RegisteredAt time.Time `json:"registered_at"`
	EditedAt     time.Time `json:"edited_at"`
}

type UpdateCategoryResponse struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Name     string    `json:"name"`
	EditedAt time.Time `json:"edited_at"`
}
