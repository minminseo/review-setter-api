package category

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Category struct {
	ID     string
	UserID string
	Name   string
}

func NewCategory(
	id string,
	userID string,
	name string,
) (*Category, error) {

	if err := validateName(name); err != nil {
		return nil, err
	}

	c := &Category{
		ID:     id,
		UserID: userID,
		Name:   name,
	}

	return c, nil
}

func ReconstructCategory(
	id string,
	userID string,
	name string,
) (*Category, error) {
	c := &Category{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
	return c, nil
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("カテゴリー名は必須です"),
	)
}

func (c *Category) Set(name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	c.Name = name
	return nil
}
