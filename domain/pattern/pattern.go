package pattern

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Pattern struct {
	PatternID    string
	UserID       string
	Name         string
	TargetWeight string
	RegisteredAt time.Time
	EditedAt     time.Time
}

func NewPattern(
	patternID string,
	userID string,
	name string,
	targetWeight string,
	registeredAt time.Time,
	editedAt time.Time,
) (*Pattern, error) {

	if err := validateName(name); err != nil {
		return nil, err
	}
	if err := validateTargetWeight(string(targetWeight)); err != nil {
		return nil, err
	}
	p := &Pattern{
		PatternID:    patternID,
		UserID:       userID,
		Name:         name,
		TargetWeight: targetWeight,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}
	return p, nil
}

func ReconstructPattern(
	patternID string,
	userID string,
	name string,
	targetWeight string,
	registeredAt time.Time,
	editedAt time.Time,
) (*Pattern, error) {
	p := &Pattern{
		PatternID:    patternID,
		UserID:       userID,
		Name:         name,
		TargetWeight: targetWeight,
		RegisteredAt: registeredAt,
		EditedAt:     editedAt,
	}
	return p, nil
}

const (
	TargetWeightHeavy  string = "heavy"
	TargetWeightNormal string = "normal"
	TargetWeightLight  string = "light"
	TargetWeightUnset  string = "unset"
)

var allowedTargetWeights = map[string]struct{}{
	TargetWeightHeavy:  {},
	TargetWeightNormal: {},
	TargetWeightLight:  {},
	TargetWeightUnset:  {},
}

func validateName(name string) error {
	return validation.Validate(
		name,
		validation.Required.Error("名前は必須です"),
	)
}
func validateTargetWeight(targetWeight string) error {
	return validation.Validate(
		targetWeight,
		validation.Required.Error("重みは必須です"),
		validation.By(func(value interface{}) error {
			trgtWght, _ := value.(string)
			if _, ok := allowedTargetWeights[trgtWght]; !ok {
				return errors.New("重みの値が不正です")
			}
			return nil
		}),
	)
}

func (p *Pattern) Set(
	name string,
	targetWeight string,
	editedAt time.Time,
) error {
	if err := validateName(name); err != nil {
		return err
	}
	if err := validateTargetWeight(string(targetWeight)); err != nil {
		return err
	}

	p.Name = name
	p.TargetWeight = targetWeight
	p.EditedAt = editedAt

	return nil
}

type PatternStep struct {
	PatternStepID string
	UserID        string
	PatternID     string
	StepNumber    int
	IntervalDays  int
}

func NewPatternStep(
	patternStepID string,
	userID string,
	patternID string,
	stepNumber int,
	intervalDays int,
) (*PatternStep, error) {
	if err := validateStepNumber(stepNumber); err != nil {
		return nil, err
	}
	if err := validateIntervalDays(intervalDays); err != nil {
		return nil, err
	}
	ps := &PatternStep{
		PatternStepID: patternStepID,
		UserID:        userID,
		PatternID:     patternID,
		StepNumber:    stepNumber,
		IntervalDays:  intervalDays,
	}

	return ps, nil
}

func ReconstructPatternStep(
	patternStepID string,
	userID string,
	patternID string,
	stepNumber int,
	intervalDays int,
) (*PatternStep, error) {
	ps := &PatternStep{
		PatternStepID: patternStepID,
		UserID:        userID,
		PatternID:     patternID,
		StepNumber:    stepNumber,
		IntervalDays:  intervalDays,
	}
	return ps, nil
}

func validateStepNumber(stepNumber int) error {
	return validation.Validate(
		stepNumber,
		validation.Required.Error("順序番号は必須です"),
		validation.Min(1).Error("順序番号の値が不正です"),
		validation.Max(32767).Error("順序番号は32768回以上は指定できません"),
	)
}
func validateIntervalDays(intervalDays int) error {
	return validation.Validate(
		intervalDays,
		validation.Required.Error("復習日間隔数は必須です"),
		validation.Min(1).Error("復習日間隔数は1以上で指定してください"),
		validation.Max(32767).Error("復習日間隔数は32768日後以上は指定できません"),
	)
}

func ValidateSteps(steps []*PatternStep) error {

	// ステップ数が0の場合はエラー
	if len(steps) == 0 {
		return errors.New("復習日間隔数は1つ以上指定してください")
	}

	// ステップ数が1つの場合は昇順チェック不要
	if len(steps) == 1 {
		return nil
	}

	// ステップ数が2つ以上の場合は昇順チェック
	prev := steps[0]
	for _, curr := range steps[1:] {
		if curr.StepNumber == prev.StepNumber {
			return errors.New("順序番号は重複してはいけません")
		}
		if curr.IntervalDays == prev.IntervalDays {
			return errors.New("復習日間隔数は重複してはいけません")
		}
		if curr.StepNumber < prev.StepNumber {
			return errors.New("順序番号は昇順で指定してください")
		}
		if curr.IntervalDays < prev.IntervalDays {
			return errors.New("復習日間隔数は昇順で指定してください")
		}

		// StepNumberが必ず前の値から+1（公差1の等差数列）になっているかチェック
		if curr.StepNumber != prev.StepNumber+1 {
			return errors.New("順序番号は1ずつ増加して指定してください")
		}

		prev = curr
	}
	return nil
}
