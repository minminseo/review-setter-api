package user

import (
	"testing"
	"time"
)

func TestNewEmailVerification(t *testing.T) {
	tests := []struct {
		name           string
		verificationID string
		userID         string
		wantErr        bool
		wantErrMsg     string
	}{
		{
			name:           "有効な認証作成（正常系）",
			verificationID: "verification1",
			userID:         "user1",
			wantErr:        false,
		},
		{
			name:           "認証IDが空（異常系）",
			verificationID: "",
			userID:         "user1",
			wantErr:        true,
			wantErrMsg:     "認証IDが空です",
		},
		{
			name:           "ユーザーIDが空（異常系）",
			verificationID: "verification2",
			userID:         "",
			wantErr:        true,
			wantErrMsg:     "ユーザーIDが空です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			before := time.Now()
			verification, code, err := NewEmailVerification(tc.verificationID, tc.userID)
			after := time.Now()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("エラーが発生することを期待しましたが、nilでした")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			// 基本フィールドの検証
			if verification.ID != tc.verificationID {
				t.Errorf("認証IDが一致しません: got %q, want %q", verification.ID, tc.verificationID)
			}
			if verification.UserID != tc.userID {
				t.Errorf("ユーザーIDが一致しません: got %q, want %q", verification.UserID, tc.userID)
			}

			// コード長の検証
			if len(code) != VerificationCodeLength {
				t.Errorf("認証コードの長さが一致しません: got %d, want %d", len(code), VerificationCodeLength)
			}

			// コードが数字のみで構成されているか
			for _, char := range code {
				if char < '0' || char > '9' {
					t.Errorf("認証コードに数字以外の文字が含まれています: %c", char)
				}
			}

			// コードハッシュの検証
			expectedHash := hashVerificationCode(code)
			if verification.CodeHash != expectedHash {
				t.Errorf("コードハッシュが一致しません: got %q, want %q", verification.CodeHash, expectedHash)
			}

			// 有効期限の検証
			beforeExpectedExpiry := before.Add(VerificationExpiry)
			afterExpectedExpiry := after.Add(VerificationExpiry)
			if verification.ExpiresAt.Before(beforeExpectedExpiry) || verification.ExpiresAt.After(afterExpectedExpiry) {
				t.Errorf("有効期限が範囲外です: got %v, expected between %v and %v", verification.ExpiresAt, beforeExpectedExpiry, afterExpectedExpiry)
			}

			// コード検証の確認
			if !verification.ValidateCode(code) {
				t.Error("生成したコードでValidateCodeがtrueを返しません")
			}
		})
	}
}

func TestEmailVerification_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "有効期限切れでない（正常系）",
			expiresAt: time.Now().Add(5 * time.Minute),
			want:      false,
		},
		{
			name:      "有効期限切れ（異常系）",
			expiresAt: time.Now().Add(-5 * time.Minute),
			want:      true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			verification := &EmailVerification{
				ID:        "verification1",
				UserID:    "user1",
				CodeHash:  "hash",
				ExpiresAt: tc.expiresAt,
			}

			got := verification.IsExpired()
			if got != tc.want {
				t.Errorf("IsExpired(): got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEmailVerification_ValidateCode(t *testing.T) {
	// Create a test verification
	verification, originalCode, err := NewEmailVerification("verification1", "user1")
	if err != nil {
		t.Fatalf("認証の生成に失敗しました: %v", err)
	}

	tests := []struct {
		name string
		code string
		want bool
	}{
		{
			name: "正しいコード（正常系）",
			code: originalCode,
			want: true,
		},
		{
			name: "間違ったコード（異常系）",
			code: "123456", // originalCodeじゃない
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := verification.ValidateCode(tc.code)
			if got != tc.want {
				t.Errorf("ValidateCode(%q): got %v, want %v", tc.code, got, tc.want)
			}
		})
	}
}

// 認証コード生成に集中させるため、6桁縛りはしなくていい
func TestGenerateVerificationCode(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "6桁（正常系）",
			length: 6,
		},
		{
			name:   "4桁（正常系）",
			length: 4,
		},
		{
			name:   "8桁（正常系）",
			length: 8,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			code, err := generateVerificationCode(tc.length)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if len(code) != tc.length {
				t.Errorf("認証コードの長さが一致しません: got %d, want %d", len(code), tc.length)
			}

			// コードが数字のみで構成されているか
			for _, char := range code {
				if char < '0' || char > '9' {
					t.Errorf("認証コードに数字以外の文字が含まれています: %c", char)
				}
			}
		})
	}
}

// hashVerificationCodeはライブラリしか使ってないロジックので、TestHashVerificationCodeはなしでいく。
