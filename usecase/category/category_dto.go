package category

import "time"

type CreateCategoryInput struct {
	UserID string
	Name   string
}

type CreateCategoryOutput struct {
	ID           string
	UserID       string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type GetCategoryOutput struct {
	ID           string
	UserID       string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type UpdateCategoryInput struct {
	ID     string
	UserID string
	Name   string
}

type UpdateCategoryOutput struct {
	ID       string
	UserID   string
	Name     string
	EditedAt time.Time
}
