package item

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	PatternDomain "github.com/minminseo/recall-setter/domain/pattern"
)

type Item struct {
	itemID       string
	userID       string
	categoryID   *string
	boxID        *string
	patternID    *string
	name         string
	detail       string
	learnedDate  time.Time
	isFinished   bool
	registeredAt time.Time
	editedAt     time.Time
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
		itemID:       itemID,
		userID:       userID,
		categoryID:   categoryID,
		boxID:        boxID,
		patternID:    patternID,
		name:         name,
		detail:       detail,
		learnedDate:  learnedDate,
		isFinished:   isFinished,
		registeredAt: registeredAt,
		editedAt:     editedAt,
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
		itemID:       itemID,
		userID:       userID,
		categoryID:   categoryID,
		boxID:        boxID,
		patternID:    patternID,
		name:         name,
		detail:       detail,
		learnedDate:  learnedDate,
		isFinished:   isFinished,
		registeredAt: registeredAt,
		editedAt:     editedAt,
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

func (i *Item) ItemID() string {
	return i.itemID
}

func (i *Item) UserID() string {
	return i.userID
}

func (i *Item) CategoryID() *string {
	return i.categoryID
}

func (i *Item) BoxID() *string {
	return i.boxID
}

func (i *Item) PatternID() *string {
	return i.patternID
}

func (i *Item) Name() string {
	return i.name
}

func (i *Item) Detail() string {
	return i.detail
}

func (i *Item) LearnedDate() time.Time {
	return i.learnedDate
}

func (i *Item) IsFinished() bool {
	return i.isFinished
}

func (i *Item) RegisteredAt() time.Time {
	return i.registeredAt
}

func (i *Item) EditedAt() time.Time {
	return i.editedAt
}

func (i *Item) UpdateItem(
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

	i.categoryID = categoryID
	i.boxID = boxID
	i.patternID = patternID
	i.name = name
	i.detail = detail
	i.learnedDate = learnedDate
	i.editedAt = editedAt

	return nil
}

type Reviewdate struct {
	reviewdateID         string
	userID               string
	categoryID           *string
	boxID                *string
	itemID               string
	stepNumber           int
	initialScheduledDate time.Time
	scheduledDate        time.Time
	isCompleted          bool
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
		reviewdateID:         reviewdateID,
		userID:               userID,
		categoryID:           categoryID,
		boxID:                boxID,
		itemID:               itemID,
		stepNumber:           stepNumber,
		initialScheduledDate: initialScheduledDate,
		scheduledDate:        scheduledDate,
		isCompleted:          isCompleted,
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
		reviewdateID:         reviewdateID,
		userID:               userID,
		categoryID:           categoryID,
		boxID:                boxID,
		itemID:               itemID,
		stepNumber:           stepNumber,
		initialScheduledDate: initialScheduledDate,
		scheduledDate:        scheduledDate,
		isCompleted:          isCompleted,
	}
	return rd, nil
}

func (r *Reviewdate) ReviewdateID() string {
	return r.reviewdateID
}

func (r *Reviewdate) UserID() string {
	return r.userID
}

func (r *Reviewdate) CategoryID() *string {
	return r.categoryID
}

func (r *Reviewdate) BoxID() *string {
	return r.boxID
}

func (r *Reviewdate) ItemID() string {
	return r.itemID
}

func (r *Reviewdate) StepNumber() int {
	return r.stepNumber
}

func (r *Reviewdate) InitialScheduledDate() time.Time {
	return r.initialScheduledDate
}

func (r *Reviewdate) ScheduledDate() time.Time {
	return r.scheduledDate
}

func (r *Reviewdate) IsCompleted() bool {
	return r.isCompleted
}

func (r *Reviewdate) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(
			&r.stepNumber,
			validation.Required.Error("ステップ番号は必須です"),
			validation.Min(1).Error("ステップ番号の値が不正です"),
			validation.Max(32767).Error("ステップは32768回以上は指定できません"),
		),
		validation.Field(
			&r.scheduledDate,
			validation.Required.Error("スケジュール日は必須です"),
			// スケジュール日フォーマットのバリデーションも書く
		),
	)
}

func (r *Reviewdate) UpdateReviewdateIDs(
	categoryID *string,
	boxID *string,
) error {
	r.categoryID = categoryID
	r.boxID = boxID
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
		diff time.Duration,
	) ([]*Reviewdate, error)
}
