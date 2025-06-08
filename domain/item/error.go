package item

import "errors"

var (
	ErrNoDiff                                     = errors.New("変更点がありません")
	ErrHasCompletedReviewDate                     = errors.New("完了済みの復習物があるため、復習パターンを変更できません")
	ErrNewScheduledDateBeforeInitialScheduledDate = errors.New("新しい復習日は初期復習日より前に設定できません")
	ErrMismatchedIDsAndSteps                      = errors.New("復習パターンのステップ数と復習日数が一致しません")
)
