package category

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Category struct {
	id     string
	userID string
	name   string
}

func NewCategory(
	id string,
	userID string,
	name string,
) (*Category, error) {
	c := &Category{
		id:     id,
		userID: userID,
		name:   name,
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Category) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(
			&c.name,
			validation.Required.Error("カテゴリー名は必須です"),
		),
	)
}
