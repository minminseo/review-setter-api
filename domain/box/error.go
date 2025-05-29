// domain/box/errors.go
package box

import "errors"

// review_itemsが存在する状態パターン変更しようとしたときのエラー
var ErrPatternConflict = errors.New("復習物がボックス内に存在するため、復習パターンを変更できません")
