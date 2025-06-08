package pattern

import "time"

type CreatePatternStepInput struct {
	StepNumber   int
	IntervalDays int
}

type CreatePatternInput struct {
	UserID       string
	Name         string
	TargetWeight string
	Steps        []CreatePatternStepInput
}

type CreatePatternStepOutput struct {
	PatternStepID string
	UserID        string
	PatternID     string
	StepNumber    int
	IntervalDays  int
}

type CreatePatternOutput struct {
	ID           string
	UserID       string
	Name         string
	TargetWeight string
	RegisteredAt time.Time
	EditedAt     time.Time
	Steps        []CreatePatternStepOutput
}

type GetPatternStepOutput struct {
	PatternStepID string
	PatternID     string
	StepNumber    int
	IntervalDays  int
}

type GetPatternOutput struct {
	PatternID    string
	UserID       string
	Name         string
	TargetWeight string
	RegisteredAt time.Time
	EditedAt     time.Time
	Steps        []GetPatternStepOutput
}

type UpdatePatternStepInput struct {
	StepID       string
	PatternID    string
	StepNumber   int
	IntervalDays int
}

type UpdatePatternInput struct {
	PatternID    string
	UserID       string
	Name         string
	TargetWeight string
	Steps        []UpdatePatternStepInput
}

type UpdatePatternStepOutput struct {
	PatternStepID string
	UserID        string
	PatternID     string
	StepNumber    int
	IntervalDays  int
}

type UpdatePatternOutput struct {
	PatternID    string
	UserID       string
	Name         string
	TargetWeight string
	RegisteredAt time.Time
	EditedAt     time.Time
	Steps        []UpdatePatternStepOutput
}
