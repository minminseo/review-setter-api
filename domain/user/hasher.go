package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct {
	key []byte
}

func NewHasher(hexKey string) (*Hasher, error) {
	key, err := hex.DecodeString(hexKey) // 秘密鍵をバイト列に変換
	if err != nil {
		return nil, err
	}
	return &Hasher{key: key}, nil
}

// メールアドレスからユーザー検索用のキーを生成
func (h *Hasher) GenerateSearchKey(email string) string {
	//ハッシュ関数と秘密鍵を設定 → macにemailを書き込む → ハッシュ値算出→16進数文字列に変換してreturn
	mac := hmac.New(sha256.New, h.key)
	mac.Write([]byte(email))
	return hex.EncodeToString(mac.Sum(nil)) // そこまでパフォは重要ではないのでnil
}
