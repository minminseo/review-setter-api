package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	VerificationCodeLength = 6
	VerificationExpiry     = 10 * time.Minute
)

type EmailVerification struct {
	ID        string
	UserID    string
	CodeHash  string
	ExpiresAt time.Time
}

// 認証情報の生成
func NewEmailVerification(verificationID string, userID string) (*EmailVerification, string, error) {
	code, err := generateVerificationCode(VerificationCodeLength)
	if err != nil {
		return nil, "", fmt.Errorf("認証コードの生成に失敗しました: %w", err)
	}

	codeHash := hashVerificationCode(code)
	expiresAt := time.Now().Add(VerificationExpiry)

	return &EmailVerification{
		ID:        verificationID,
		UserID:    userID,
		CodeHash:  codeHash,
		ExpiresAt: expiresAt,
	}, code, nil
}

// 有効期限が切れているか確認
func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.ExpiresAt)
}

// 認証コード検証
func (ev *EmailVerification) ValidateCode(code string) bool {
	return ev.CodeHash == hashVerificationCode(code)
}

// 認証コード生成
func generateVerificationCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := ""
	for _, b := range bytes {
		code += fmt.Sprintf("%d", b%10)
	}
	return code, nil
}

// 認証コードハッシュ化
func hashVerificationCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}
