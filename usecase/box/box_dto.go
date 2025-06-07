package box

import "time"

type CreateBoxInput struct {
	UserID     string
	CategoryID string
	PatternID  string
	Name       string
}

type CreateBoxOutput struct {
	ID           string
	UserID       string
	CategoryID   string
	PatternID    string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type GetBoxOutput struct {
	ID           string
	UserID       string
	CategoryID   string
	PatternID    string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

type UpdateBoxInput struct {
	ID         string
	UserID     string
	CategoryID string
	PatternID  string
	Name       string
}

type UpdateBoxOutput struct {
	ID         string
	UserID     string
	CategoryID string
	PatternID  string
	Name       string
	EditedAt   time.Time
}
