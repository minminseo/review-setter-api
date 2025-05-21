package pattern

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type targetWeightName string

type Pattern struct {
	id           string
	userID       string
	name         string
	targetWeight targetWeightName
}

func NewPattern(
	id string,
	userID string,
	name string,
	targetWeight targetWeightName,
) (*Pattern, error) {
	p := &Pattern{
		id:           id,
		userID:       userID,
		name:         name,
		targetWeight: targetWeight,
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

const (
	TargetWeightHeavy  targetWeightName = "heavy"
	TargetWeightNormal targetWeightName = "normal"
	TargetWeightLight  targetWeightName = "light"
	TargetWeightUnset  targetWeightName = "unset"
)

var allowedTargetWeights = map[targetWeightName]struct{}{
	TargetWeightHeavy:  {},
	TargetWeightNormal: {},
	TargetWeightLight:  {},
	TargetWeightUnset:  {},
}

func (p *Pattern) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(
			&p.name,
			validation.Required.Error("名前は必須です"),
		),
		validation.Field(
			&p.targetWeight,
			validation.Required.Error("重みは必須です"),
			validation.By(func(value interface{}) error {
				trgtWght, _ := value.(targetWeightName)
				if _, ok := allowedTargetWeights[trgtWght]; !ok {
					return errors.New("重みの値が不正です")
				}
				return nil
			}),
		),
	)
}

type PatternStep struct {
	id           string
	patternID    string
	stepNumber   int
	intervalDays int
}

func NewPatternStep(
	id string,
	patternID string,
	stepNumber int,
	intervalDays int,
) (*PatternStep, error) {
	ps := &PatternStep{
		id:           id,
		patternID:    patternID,
		stepNumber:   stepNumber,
		intervalDays: intervalDays,
	}
	if err := ps.Validate(); err != nil {
		return nil, err
	}
	return ps, nil
}

func (ps *PatternStep) Validate() error {
	return validation.ValidateStruct(ps,
		validation.Field(
			&ps.stepNumber,
			validation.Required.Error("ステップ番号は必須です"),
			validation.Min(1).Error("ステップ番号の値が不正です"),
			validation.Max(32767).Error("ステップは32768回以上は指定できません"),
		),
		validation.Field(
			&ps.intervalDays,
			validation.Required.Error("間隔日数は必須です"), //　名前が微妙
			validation.Min(1).Error("間隔日数は1以上で指定してください"),
			validation.Max(32767).Error("間隔日数は32768日後以上は指定できません"),
		),
	)
}

func validateSteps(steps []*PatternStep) error {

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
		if curr.stepNumber == prev.stepNumber {
			return errors.New("順序番号は重複してはいけません")
		}
		if curr.intervalDays == prev.intervalDays {
			return errors.New("復習日間隔数は重複してはいけません")
		}
		if curr.stepNumber < prev.stepNumber {
			return errors.New("順序番号は昇順で指定してください")
		}
		if curr.intervalDays < prev.intervalDays {
			return errors.New("復習日間隔数は昇順で指定してください")
		}

		// stepNumberが必ず前の値から+1（公差1の等差数列）になっているかチェック
		if curr.stepNumber != prev.stepNumber+1 {
			return errors.New("順序番号は1ずつ増加して指定してください")
		}

		prev = curr
	}
	return nil
}
