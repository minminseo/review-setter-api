package pattern

type CreatePatternRequest struct {
	Name         string                   `json:"name"`
	TargetWeight string                   `json:"target_weight"`
	Steps        []CreatePatternStepField `json:"steps"`
}
type CreatePatternStepField struct {
	StepNumber   int `json:"step_number"`
	IntervalDays int `json:"interval_days"`
}

type UpdatePatternRequest struct {
	Name         string                   `json:"name"`
	TargetWeight string                   `json:"target_weight"`
	Steps        []UpdatePatternStepField `json:"steps"`
}
type UpdatePatternStepField struct {
	StepID       string `json:"step_id"`
	StepNumber   int    `json:"step_number"`
	IntervalDays int    `json:"interval_days"`
}
