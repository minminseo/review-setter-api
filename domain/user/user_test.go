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
	// テスト用の暗号化サービスとハッシュサービスを作成
	cryptoService, err := NewCryptoService("82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d")
	if err != nil {
		t.Fatalf("failed to create crypto service: %v", err)
	}

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
			name:       "有効な入力（正常系）",
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
			name:       "無効なメールアドレス（異常系）",
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
			name:       "短いパスワード（異常系）",
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
			name:       "サポートされていないタイムゾーン（異常系）",
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
			name:       "サポートされていないテーマカラー（異常系）",
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
			name:       "サポートされていない言語（異常系）",
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
			t.Parallel()

			searchKey := ""
			u, err := NewUser(tc.id, tc.email, tc.password, tc.timezone, tc.themeColor, tc.language, cryptoService, searchKey)
			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待しましたが、nilでした")
				}
				if err.Error() != tc.errMsg {
					t.Errorf("エラーメッセージが一致しません: got %q, want %q", err.Error(), tc.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			// フィールドのセットを確認
			if u.ID != tc.id {
				t.Errorf("ユーザーIDが一致しません: got %q, want %q", u.ID, tc.id)
			}

			// メールアドレスは暗号化されているので復号化して確認
			decryptedEmail, err := u.GetEmail(cryptoService)
			if err != nil {
				t.Fatalf("メールアドレスの復号に失敗しました: %v", err)
			}
			if decryptedEmail != tc.email {
				t.Errorf("メールアドレスが一致しません: got %q, want %q", decryptedEmail, tc.email)
			}

			if u.Timezone != tc.timezone {
				t.Errorf("タイムゾーンが一致しません: got %q, want %q", u.Timezone, tc.timezone)
			}
			if u.ThemeColor != tc.themeColor {
				t.Errorf("テーマカラーが一致しません: got %q, want %q", u.ThemeColor, tc.themeColor)
			}
			if u.Language != tc.language {
				t.Errorf("言語タグが一致しません: got %q, want %q", u.Language, tc.language)
			}

			// 初期状態では未認証であることを確認
			if u.IsVerified() {
				t.Error("新規ユーザーは未認証であるべきです")
			}

			// パスワードがハッシュ化されていること & IsValidPasswordで検証できること
			if err := u.IsValidPassword(tc.password); err != nil {
				t.Errorf("正しいパスワードでIsValidPasswordがエラーを返しました: %v", err)
			}
			if err := u.IsValidPassword("wrongpass"); err == nil {
				t.Errorf("間違ったパスワードでIsValidPasswordがエラーを返しませんでした")
			}
		})
	}
}
