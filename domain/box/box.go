package box

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Box struct {
	id           string
	userID       string
	categoryID   string
	patternID    string
	name         string
	registeredAt time.Time
	editedAt     time.Time
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
		id:           id,
		userID:       userID,
		categoryID:   categoryID,
		patternID:    patternID,
		name:         name,
		registeredAt: registeredAt,
		editedAt:     editedAt,
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
		id:           id,
		userID:       userID,
		categoryID:   categoryID,
		patternID:    patternID,
		name:         name,
		registeredAt: registeredAt,
		editedAt:     editedAt,
	}
	return b, nil
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("カテゴリー名は必須です"),
	)
}

func (b *Box) ID() string {
	return b.id
}

func (b *Box) UserID() string {
	return b.userID
}

func (b *Box) CategoryID() string {
	return b.categoryID
}

func (b *Box) PatternID() string {
	return b.patternID
}

func (b *Box) Name() string {
	return b.name
}

func (b *Box) RegisteredAt() time.Time {
	return b.registeredAt
}

func (b *Box) EditedAt() time.Time {
	return b.editedAt
}

func (b *Box) UpdateBox(
	patternID string,
	name string,
	editedAt time.Time,
) (bool, error) {
	var isSamePattern bool
	if b.patternID == patternID {
		isSamePattern = true
	} else {
		isSamePattern = false
	}

	if err := validateName(name); err != nil {
		return isSamePattern, err
	}

	b.patternID = patternID
	b.name = name
	b.editedAt = editedAt

	return isSamePattern, nil
}
