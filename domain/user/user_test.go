package user

// ozzo-validationのis.Emailの条件
/*
"plainaddress",             // @ がない
"@missingusername.com",     // ユーザー名がない
"user@.com",                // ドメインがピリオドで始まる
"user@com",                 // ドメインにピリオドがない
"user..name@example.com",   // ローカル部に連続ドット
".user@example.com",        // ローカル部がピリオドで始まる
"user@example..com",        // ドメイン部に連続ドット
"user@-example.com",        // ドメインラベルがハイフンで始まる
"user@exam_ple.com",        // ドメインに許可されない文字（アンダースコア）
7文字以上、254文字以下の長さがOK
*/

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		email      string
		password   string
		timezone   string
		themeColor string
		language   string
		wantErr    bool
		errMsg     string
	}{
		// 正常系
		{
			name:       "valid input",
			id:         "user1",
			email:      "test@example.com",
			password:   "secret123",
			timezone:   TimeZoneTokyo,
			themeColor: ThemeColorDark,
			language:   LanguageJa,
			wantErr:    false,
		},

		// 異常系
		{
			name:       "invalid email",
			id:         "user2",
			email:      "invalid-email",
			password:   "secret123",
			timezone:   TimeZoneTokyo,
			themeColor: ThemeColorDark,
			language:   LanguageJa,
			wantErr:    true,
			errMsg:     "メールアドレスを入力して下さい",
		},

		// 異常系
		{
			name:       "short password",
			id:         "user3",
			email:      "abcdefg@example.com",
			password:   "123",
			timezone:   TimeZoneTokyo,
			themeColor: ThemeColorDark,
			language:   LanguageJa,
			wantErr:    true,
			errMsg:     "パスワードは6文字以上です",
		},

		// 異常系
		{
			name:       "unsupported timezone",
			id:         "user4",
			email:      "abcdefg@example.com",
			password:   "secret123",
			timezone:   "Invalid/Zone",
			themeColor: ThemeColorDark,
			language:   LanguageJa,
			wantErr:    true,
			errMsg:     "タイムゾーンの値が不正です",
		},
		{
			name:       "unsupported theme color",
			id:         "user5",
			email:      "abcdefg@example.com",
			password:   "secret123",
			timezone:   TimeZoneTokyo,
			themeColor: "blue",
			language:   LanguageJa,
			wantErr:    true,
			errMsg:     "テーマカラーは'dark'または'light'で指定してください",
		},

		// 異常系
		{
			name:       "unsupported language",
			id:         "user6",
			email:      "abcdefg@example.com",
			password:   "secret123",
			timezone:   TimeZoneTokyo,
			themeColor: ThemeColorDark,
			language:   "test",
			wantErr:    true,
			errMsg:     "言語タグの値が不正です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			u, err := NewUser(tc.id, tc.email, tc.password, tc.timezone, tc.themeColor, tc.language)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("unexpected error message: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// フィールドのセットを確認
			if u.ID != tc.id {
				t.Errorf("ID: got %q, want %q", u.ID, tc.id)
			}
			if u.Email != tc.email {
				t.Errorf("Email: got %q, want %q", u.Email, tc.email)
			}
			if u.Timezone != tc.timezone {
				t.Errorf("Timezone: got %q, want %q", u.Timezone, tc.timezone)
			}
			if u.ThemeColor != tc.themeColor {
				t.Errorf("ThemeColor: got %q, want %q", u.ThemeColor, tc.themeColor)
			}
			if u.Language != tc.language {
				t.Errorf("Language: got %q, want %q", u.Language, tc.language)
			}
			// パスワードがハッシュ化されていること & IsValidPasswordで検証できること
			if err := u.IsValidPassword(tc.password); err != nil {
				t.Errorf("IsValidPassword returned error for correct password: %v", err)
			}
			if err := u.IsValidPassword("wrongpass"); err == nil {
				t.Errorf("IsValidPassword did not return error for wrong password")
			}
		})
	}
}
