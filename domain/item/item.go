package item

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Item struct {
	id          string
	userID      string
	categoryID  string
	boxID       string
	patternID   string
	name        string
	detail      string
	learnedDate string
	isFinished  bool
}

func NewItem(
	id string,
	userID string,
	categoryID string,
	boxID string,
	patternID string,
	name string,
	detail string,
	learnedDate string,
	isFinished bool,
) (*Item, error) {
	i := &Item{
		id:          id,
		userID:      userID,
		categoryID:  categoryID,
		boxID:       boxID,
		patternID:   patternID,
		name:        name,
		detail:      detail,
		learnedDate: learnedDate,
		isFinished:  isFinished,
	}
	if err := i.Validate(); err != nil {
		return nil, err
	}
	return i, nil
}
func (i *Item) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(
			&i.name,
			validation.Required.Error("アイテム名は必須です"),
		),
		validation.Field(
			&i.learnedDate,
			validation.Required.Error("学習日は必須です"),
		),
		validation.Field(
			&i.isFinished,
			validation.Required.Error("完了フラグは必須です"),
		),
	)
}
