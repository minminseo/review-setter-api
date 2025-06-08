package pattern

import "errors"

var (
	ErrNoDiff               = errors.New("変更点がありません")
	ErrPatternNotFound      = errors.New("復習パターンが存在しません")
	ErrPatternRelatedToItem = errors.New("この復習パターンは復習物に紐づいているため削除できません")
)
