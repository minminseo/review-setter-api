package box

import validation "github.com/go-ozzo/ozzo-validation/v4"

type box struct {
	id         string
	userID     string
	categoryID string
	patternID  string
	name       string
}

func NewBox(
	id string,
	userID string,
	categoryID string,
	patternID string,
	name string,
) (*box, error) {
	b := &box{
		id:         id,
		userID:     userID,
		categoryID: categoryID,
		patternID:  patternID,
		name:       name,
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *box) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(
			&b.name,
			validation.Required.Error("ボックス名は必須です"),
		),
	)
}
