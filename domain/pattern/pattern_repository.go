package pattern

type PatternInfo struct {
	PatternID     string
	PatternStepID string
	Name          string
	TargetWeight  string
	StepNumber    int
	IntervalDays  int
}

type PatternRepository interface {
	Create(pattern *Pattern) (*Pattern, error)
	Update(pattern *Pattern, userID string) (*Pattern, error)
	Delete(patternID string, userID string) error
	GetAllByUserID(userID string) ([]*PatternInfo, error)

	// ボックス一覧取得→ボックス毎にループ処理（Patternを取得→PatternStepたちを取得）
	FindPatternByPatternID(patternID string) (*Pattern, error)
	GetAllPatternStepByPatternID(patternID string) ([]*PatternStep, error)
}
