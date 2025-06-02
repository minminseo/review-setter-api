package pattern

import "errors"

var (
	ErrNoDiff          = errors.New("変更点がありません")
	ErrPatternNotFound = errors.New("復習パターンが存在しません")
)
