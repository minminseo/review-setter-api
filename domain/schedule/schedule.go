package schedule

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Schedule struct {
	id            string
	itemID        string
	stepNumber    int
	scheduledDate string
	isCompleted   bool
}

func NewSchedule(
	id string,
	itemID string,
	stepNumber int,
	scheduledDate string,
	isCompleted bool,
) (*Schedule, error) {
	s := &Schedule{
		id:            id,
		itemID:        itemID,
		stepNumber:    stepNumber,
		scheduledDate: scheduledDate,
		isCompleted:   isCompleted,
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Schedule) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(
			&s.stepNumber,
			validation.Required.Error("ステップ番号は必須です"),
			validation.Min(1).Error("ステップ番号の値が不正です"),
			validation.Max(32767).Error("ステップは32768回以上は指定できません"),
		),
		validation.Field(
			&s.scheduledDate,
			validation.Required.Error("スケジュール日は必須です"),
			// スケジュール日フォーマットのバリデーションも書く
		),
	)
}
