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
	id        string
	userID    string
	codeHash  string
	expiresAt time.Time
}

// 認証情報の生成
func NewEmailVerification(verificationID string, userID string) (*EmailVerification, string, error) {

	if verificationID == "" {
		return nil, "", fmt.Errorf("認証IDが空です")
	}
	if userID == "" {
		return nil, "", fmt.Errorf("ユーザーIDが空です")
	}

	code, err := generateVerificationCode(VerificationCodeLength)
	if err != nil {
		return nil, "", fmt.Errorf("認証コードの生成に失敗しました: %w", err)
	}

	codeHash := hashVerificationCode(code)
	expiresAt := time.Now().Add(VerificationExpiry)

	return &EmailVerification{
		id:        verificationID,
		userID:    userID,
		codeHash:  codeHash,
		expiresAt: expiresAt,
	}, code, nil
}

// リポジトリからの復元用
func ReconstructEmailVerification(
	id string,
	userID string,
	codeHash string,
	expiresAt time.Time,
) (*EmailVerification, error) {
	if id == "" {
		return nil, fmt.Errorf("認証IDが空です")
	}
	if userID == "" {
		return nil, fmt.Errorf("ユーザーIDが空です")
	}

	return &EmailVerification{
		id:        id,
		userID:    userID,
		codeHash:  codeHash,
		expiresAt: expiresAt,
	}, nil
}

func (ev *EmailVerification) ID() string {
	return ev.id
}

func (ev *EmailVerification) UserID() string {
	return ev.userID
}

func (ev *EmailVerification) CodeHash() string {
	return ev.codeHash
}

func (ev *EmailVerification) ExpiresAt() time.Time {
	return ev.expiresAt
}

func (ev *EmailVerification) IsExpired() bool {
	return time.Now().After(ev.expiresAt)
}

// 認証コード検証
func (ev *EmailVerification) ValidateCode(code string) bool {
	return ev.codeHash == hashVerificationCode(code)
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

// テスト用にハッシュ関数をエクスポート
func HashVerificationCodeForTest(code string) string {
	return hashVerificationCode(code)
}
