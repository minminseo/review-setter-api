package category

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Category struct {
	ID           string
	UserID       string
	Name         string
	RegisteredAt time.Time
	EditedAt     time.Time
}

func NewCategory(
	id string,
	userID string,
	name string,
	registeredAt time.Time,
	editedAt time.Time,
) (*Category, error) {

	if err := validateName(name); err != nil {
		return nil, err
	}

	c := &Category{
		ID:           id,
		UserID:       userID,
		Name:         name,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}

	return c, nil
}

func ReconstructCategory(
	id string,
	userID string,
	name string,
	registeredAt time.Time,
	editedAt time.Time,
) (*Category, error) {
	c := &Category{
		ID:           id,
		UserID:       userID,
		Name:         name,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}
	return c, nil
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("カテゴリー名は必須です"),
	)
}

func (c *Category) Set(name string, editedAt time.Time) error {
	if err := validateName(name); err != nil {
		return err
	}

	c.Name = name
	c.EditedAt = editedAt
	return nil
}
