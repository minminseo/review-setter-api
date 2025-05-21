/*

CREATE TABLE review_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES users(id) ON DELETE SET NULL,
    box_id UUID REFERENCES users(id) ON DELETE SET NULL,
    pattern_id UUID REFERENCES review_patterns(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    detail TEXT,
    learned_date DATE NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)
*/

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
	isCompleted bool
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
	isCompleted bool,
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
		isCompleted: isCompleted,
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
			&i.isCompleted,
			validation.Required.Error("完了フラグは必須です"),
		),
	)
}
