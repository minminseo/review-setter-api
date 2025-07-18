package item

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
)

type Item struct {
	ItemID       string
	UserID       string
	CategoryID   *string
	BoxID        *string
	PatternID    *string
	Name         string
	Detail       string
	LearnedDate  time.Time
	IsFinished   bool
	RegisteredAt time.Time
	EditedAt     time.Time
}

func NewItem(
	itemID string,
	userID string,
	categoryID *string,
	boxID *string,
	patternID *string,
	name string,
	detail string,
	learnedDate time.Time,
	isFinished bool,
	registeredAt time.Time,
	editedAt time.Time,
) (*Item, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}
	if err := validateLearnedDate(learnedDate); err != nil {
		return nil, err
	}

	i := &Item{
		ItemID:       itemID,
		UserID:       userID,
		CategoryID:   categoryID,
		BoxID:        boxID,
		PatternID:    patternID,
		Name:         name,
		Detail:       detail,
		LearnedDate:  learnedDate,
		IsFinished:   isFinished,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}
	return i, nil
}

func ReconstructItem(
	itemID string,
	userID string,
	categoryID *string,
	boxID *string,
	patternID *string,
	name string,
	detail string,
	learnedDate time.Time,
	isFinished bool,
	registeredAt time.Time,
	editedAt time.Time,
) (*Item, error) {
	i := &Item{
		ItemID:       itemID,
		UserID:       userID,
		CategoryID:   categoryID,
		BoxID:        boxID,
		PatternID:    patternID,
		Name:         name,
		Detail:       detail,
		LearnedDate:  learnedDate,
		IsFinished:   isFinished,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}
	return i, nil
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("復習物名は必須です"),
	)
}

func validateLearnedDate(learnedDate time.Time) error {
	return validation.Validate(
		learnedDate,
		validation.Required.Error("学習日は必須です"),
	)
}

// TODO: bool値用のバリデーション

func (i *Item) Set(
	categoryID *string,
	boxID *string,
	patternID *string,
	name string,
	detail string,
	learnedDate time.Time,
	editedAt time.Time,
) error {
	if err := validateName(name); err != nil {
		return err
	}
	if err := validateLearnedDate(learnedDate); err != nil {
		return err
	}

	i.CategoryID = categoryID
	i.BoxID = boxID
	i.PatternID = patternID
	i.Name = name
	i.Detail = detail
	i.LearnedDate = learnedDate
	i.EditedAt = editedAt

	return nil
}

type Reviewdate struct {
	ReviewdateID         string
	UserID               string
	CategoryID           *string
	BoxID                *string
	ItemID               string
	StepNumber           int
	InitialScheduledDate time.Time
	ScheduledDate        time.Time
	IsCompleted          bool
}

func NewReviewdate(
	reviewdateID string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	stepNumber int,
	initialScheduledDate time.Time,
	scheduledDate time.Time,
	isCompleted bool,
) (*Reviewdate, error) {
	s := &Reviewdate{
		ReviewdateID:         reviewdateID,
		UserID:               userID,
		CategoryID:           categoryID,
		BoxID:                boxID,
		ItemID:               itemID,
		StepNumber:           stepNumber,
		InitialScheduledDate: initialScheduledDate,
		ScheduledDate:        scheduledDate,
		IsCompleted:          isCompleted,
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func ReconstructReviewdate(
	reviewdateID string,
	userID string,
	categoryID *string,
	boxID *string,
	itemID string,
	stepNumber int,
	initialScheduledDate time.Time,
	scheduledDate time.Time,
	isCompleted bool,
) (*Reviewdate, error) {
	rd := &Reviewdate{
		ReviewdateID:         reviewdateID,
		UserID:               userID,
		CategoryID:           categoryID,
		BoxID:                boxID,
		ItemID:               itemID,
		StepNumber:           stepNumber,
		InitialScheduledDate: initialScheduledDate,
		ScheduledDate:        scheduledDate,
		IsCompleted:          isCompleted,
	}
	return rd, nil
}

func (s *Reviewdate) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(
			&s.StepNumber,
			validation.Required.Error("ステップ番号は必須です"),
			validation.Min(1).Error("ステップ番号の値が不正です"),
			validation.Max(32767).Error("ステップは32768回以上は指定できません"),
		),
		validation.Field(
			&s.ScheduledDate,
			validation.Required.Error("スケジュール日は必須です"),
			// スケジュール日フォーマットのバリデーションも書く
		),
	)
}

func (s *Reviewdate) SetOnlyIDs(
	categoryID *string,
	boxID *string,
) error {
	s.CategoryID = categoryID
	s.BoxID = boxID
	return nil
}

type IScheduler interface {
	FormatWithOverdueMarkedCompleted(
		targetPatternSteps []*PatternDomain.PatternStep,
		userID string,
		categoryID *string,
		boxID *string,
		itemID string,
		parsedLearnedDate time.Time,
		parsedToday time.Time,
	) ([]*Reviewdate, bool, error)

	FormatWithOverdueMarkedInCompleted(
		targetPatternSteps []*PatternDomain.PatternStep,
		userID string,
		categoryID *string,
		boxID *string,
		itemID string,
		parsedLearnedDate time.Time,
		parsedToday time.Time,
	) ([]*Reviewdate, error)

	FormatWithOverdueMarkedCompletedWithIDs(
		targetPatternSteps []*PatternDomain.PatternStep,
		reviewDateIDs []string,
		userID string,
		categoryID *string,
		boxID *string,
		itemID string,
		parsedLearnedDate time.Time,
		parsedToday time.Time,
	) ([]*Reviewdate, bool, error)

	FormatWithOverdueMarkedInCompletedWithIDs(
		targetPatternSteps []*PatternDomain.PatternStep,
		reviewDateIDs []string,
		userID string,
		categoryID *string,
		boxID *string,
		itemID string,
		parsedLearnedDate time.Time,
		parsedToday time.Time,
	) ([]*Reviewdate, error)

	FormatWithOverdueMarkedInCompletedWithIDsForBackReviewDates(
		targetPatternSteps []*PatternDomain.PatternStep,
		reviewDateIDs []string,
		userID string,
		categoryID *string,
		boxID *string,
		itemID string,
		parsedLearnedDate time.Time,
		parsedToday time.Time,
	) ([]*Reviewdate, error)
}
