package box

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Box struct {
	ID           string
	UserID       string
	CategoryID   string
	PatternID    string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

func NewBox(
	id string,
	userID string,
	categoryID string,
	patternID string,
	name string,
	registeredAt time.Time,
	editedAt time.Time,
) (*Box, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	b := &Box{
		ID:           id,
		UserID:       userID,
		CategoryID:   categoryID,
		PatternID:    patternID,
		Name:         name,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}

	return b, nil
}

func ReconstructBox(
	id string,
	userID string,
	categoryID string,
	patternID string,
	name string,
	registeredAt time.Time,
	editedAt time.Time,
) (*Box, error) {
	b := &Box{
		ID:           id,
		UserID:       userID,
		CategoryID:   categoryID,
		PatternID:    patternID,
		Name:         name,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}
	return b, nil
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("カテゴリー名は必須です"),
	)
}

func (b *Box) Set(
	patternID string,
	name string,
	editedAt time.Time,
) (bool, error) {
	var isSamePattern bool
	if b.PatternID == patternID {
		isSamePattern = true
	} else {
		isSamePattern = false
	}

	if err := validateName(name); err != nil {
		return isSamePattern, err
	}

	b.PatternID = patternID
	b.Name = name
	b.EditedAt = editedAt

	return isSamePattern, nil
}
