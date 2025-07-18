package user

import (
	"testing"
)

func TestNewHasher(t *testing.T) {
	tests := []struct {
		name    string
		hexKey  string
		wantErr bool
	}{
		{
			name:    "有効な16進数キー（正常系）",
			hexKey:  "edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce",
			wantErr: false,
		},
		{
			name:    "短い16進数キー（正常系）",
			hexKey:  "edc1a05c",
			wantErr: false,
		},
		{
			name:    "長い16進数キー（正常系）",
			hexKey:  "edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce12345678",
			wantErr: false,
		},
		{
			name:    "無効な16進数文字列（異常系）",
			hexKey:  "invalid_hex_string",
			wantErr: true,
		},
		{
			name:    "空文字キー（正常系）",
			hexKey:  "",
			wantErr: false,
		},
		{
			name:    "大文字小文字混在16進数キー（正常系）",
			hexKey:  "EdC1a05C210fA9b25FfA8DbbD5682BdeB393Dc77",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			hasher, err := NewHasher(tc.hexKey)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("エラーが発生することを期待しましたが、nilでした")
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}

			if hasher == nil {
				t.Fatal("hasherがnilです")
			}
		})
	}
}

func TestHasher_GenerateSearchKey(t *testing.T) {
	hasher, err := NewHasher("edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce")
	if err != nil {
		t.Fatalf("ハッシャーの生成に失敗しました: %v", err)
	}

	tests := []struct {
		name  string
		email string
	}{
		{
			name:  "シンプルなメールアドレス（正常系）",
			email: "test@example.com",
		},
		{
			name:  "複雑なメールアドレス（正常系）",
			email: "user.name+tag@sub.domain.example.com",
		},
		{
			name:  "大文字メールアドレス（正常系）",
			email: "TEST@EXAMPLE.COM",
		},
		{
			name:  "日本語メールアドレス（正常系）",
			email: "テスト@example.com",
		},
		{
			name:  "空文字メールアドレス（正常系）",
			email: "",
		},
		{
			name:  "特殊文字メールアドレス（正常系）",
			email: "user+test@example.com",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			searchKey := hasher.GenerateSearchKey(tc.email)

			// 検索キーが空でないことを確認
			if searchKey == "" {
				t.Fatal("検索キーが空です")
			}

			// 16進数文字列であることを確認（SHA256なので64文字）
			if len(searchKey) != 64 {
				t.Errorf("検索キーの長さ: got %d, want 64", len(searchKey))
			}

			// 有効な16進数文字列であることを確認
			if !isValidHex(searchKey) {
				t.Errorf("検索キーが16進数文字列ではありません: %s", searchKey)
			}

			// 同じメールアドレスは同じ検索キーを生成することを確認
			searchKey2 := hasher.GenerateSearchKey(tc.email)
			if searchKey != searchKey2 {
				t.Errorf("同じメールアドレスで異なる検索キーが生成されました: got %s and %s", searchKey, searchKey2)
			}
		})
	}
}

func TestHasher_GenerateSearchKey_Consistency(t *testing.T) {
	hasher, err := NewHasher("edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce")
	if err != nil {
		t.Fatalf("ハッシャーの生成に失敗しました: %v", err)
	}

	email := "test@example.com"

	// 複数回実行して常に同じ結果が得られることを確認
	results := make([]string, 10)
	for i := 0; i < 10; i++ {
		results[i] = hasher.GenerateSearchKey(email)
	}

	first := results[0]
	for i, result := range results {
		if result != first {
			t.Errorf("%d回目: got %s, want %s", i, result, first)
		}
	}
}

func TestHasher_GenerateSearchKey_Uniqueness(t *testing.T) {
	hasher, err := NewHasher("edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce")
	if err != nil {
		t.Fatalf("ハッシャーの生成に失敗しました: %v", err)
	}

	emails := []string{
		"test1@example.com",
		"test2@example.com",
		"test@example1.com",
		"test@example2.com",
		"TEST@EXAMPLE.COM",
		"test@EXAMPLE.COM",
	}

	searchKeys := make(map[string]string)

	for _, email := range emails {
		searchKey := hasher.GenerateSearchKey(email)

		// 異なるメールアドレスは異なる検索キーを生成することを確認
		for existingEmail, existingKey := range searchKeys {
			if email != existingEmail && searchKey == existingKey {
				t.Errorf("異なるメールアドレスで同じ検索キーが生成されました: %s と %s が %s",
					email, existingEmail, searchKey)
			}
		}

		searchKeys[email] = searchKey
	}
}

func TestHasher_GenerateSearchKey_DifferentKeys(t *testing.T) {
	// 異なる秘密鍵を持つHasherが異なる検索キーを生成することを確認
	hasher1, err := NewHasher("edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce")
	if err != nil {
		t.Fatalf("1つ目のハッシャーの生成に失敗しました: %v", err)
	}

	hasher2, err := NewHasher("123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0")
	if err != nil {
		t.Fatalf("2つ目のハッシャーの生成に失敗しました: %v", err)
	}

	email := "test@example.com"

	searchKey1 := hasher1.GenerateSearchKey(email)
	searchKey2 := hasher2.GenerateSearchKey(email)

	if searchKey1 == searchKey2 {
		t.Error("異なるハッシャーキーで同じ検索キーが生成されました")
	}
}

func TestHasher_GenerateSearchKey_SpecialCases(t *testing.T) {
	hasher, err := NewHasher("edc1a05c210fa9b25ffa8dbbd5682bdeb393dc77cbf4694599ada5a32fb48dce")
	if err != nil {
		t.Fatalf("ハッシャーの生成に失敗しました: %v", err)
	}

	tests := []struct {
		name        string
		email       string
		expectDiff  string
		description string
	}{
		{
			name:        "大文字小文字の違い（異なるキー）",
			email:       "test@example.com",
			expectDiff:  "TEST@EXAMPLE.COM",
			description: "大文字小文字の違いで異なる検索キーになること",
		},
		{
			name:        "サブドメインの違い（異なるキー）",
			email:       "test@example.com",
			expectDiff:  "test@sub.example.com",
			description: "サブドメインが異なると異なる検索キーになること",
		},
		{
			name:        "ローカル部の違い（異なるキー）",
			email:       "test@example.com",
			expectDiff:  "test1@example.com",
			description: "ローカル部が異なると異なる検索キーになること",
		},
		{
			name:        "プラスアドレッシング（異なるキー）",
			email:       "test@example.com",
			expectDiff:  "test+tag@example.com",
			description: "プラスアドレッシングで異なる検索キーになること",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			key1 := hasher.GenerateSearchKey(tc.email)
			key2 := hasher.GenerateSearchKey(tc.expectDiff)

			if key1 == key2 {
				t.Errorf("%s: メールアドレス%qと%qで同じ検索キー%qが生成されました",
					tc.description, tc.email, tc.expectDiff, key1)
			}
		})
	}
}

// ヘルパー関数: 有効な16進数文字列かチェック
func isValidHex(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}
