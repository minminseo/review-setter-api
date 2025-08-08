package category

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Category struct {
	id           string
	userID       string
	name         string
	registeredAt time.Time
	editedAt     time.Time
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
		id:           id,
		userID:       userID,
		name:         name,
		registeredAt: registeredAt,
		editedAt:     editedAt,
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
		id:           id,
		userID:       userID,
		name:         name,
		registeredAt: registeredAt,
		editedAt:     editedAt,
	}
	return c, nil
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("カテゴリー名は必須です"),
	)
}

func (c *Category) ID() string {
	return c.id
}

func (c *Category) UserID() string {
	return c.userID
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) RegisteredAt() time.Time {
	return c.registeredAt
}

func (c *Category) EditedAt() time.Time {
	return c.editedAt
}

func (c *Category) UpdateCategory(name string, editedAt time.Time) error {
	if err := validateName(name); err != nil {
		return err
	}

	c.name = name
	c.editedAt = editedAt
	return nil
}
