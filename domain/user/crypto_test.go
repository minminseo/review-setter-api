package user

import (
	"strings"
	"testing"
)

const testEmail = "test@example.com"

func TestNewCryptoService(t *testing.T) {
	tests := []struct {
		name    string
		hexKey  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "32バイトの有効な16進数キー（正常系）",
			hexKey:  "82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d",
			wantErr: false,
		},
		{
			name:    "無効な16進数文字列（異常系）",
			hexKey:  "invalid_hex_string",
			wantErr: true,
			errMsg:  "16進数の鍵のデコードに失敗しました",
		},
		{
			name:    "短すぎるキー（異常系）",
			hexKey:  "82a4d4469e7ab4916e72b05c",
			wantErr: true,
			errMsg:  "鍵のサイズが不正です: 32バイトである必要があります",
		},
		{
			name:    "長すぎるキー（異常系）",
			hexKey:  "82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d12345678",
			wantErr: true,
			errMsg:  "鍵のサイズが不正です: 32バイトである必要があります",
		},
		{
			name:    "空文字キー（異常系）",
			hexKey:  "",
			wantErr: true,
			errMsg:  "鍵のサイズが不正です: 32バイトである必要があります",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service, err := NewCryptoService(tc.hexKey)

			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待しましたが、nilでした")
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("エラーメッセージが一致しません: got %q, want to contain %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("暗号化サービスの作成に失敗しました: %v", err)
			}

			if service == nil {
				t.Fatal("serviceがnilです")
			}

			if len(service.key) != 32 {
				t.Errorf("鍵の長さが一致しません: got %d, want 32", len(service.key))
			}
		})
	}
}

func TestCryptoService_Encrypt(t *testing.T) {
	service, err := NewCryptoService("82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d")
	if err != nil {
		t.Fatalf("暗号化サービスの作成に失敗しました: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "有効なメールアドレス（正常系）",
			plaintext: testEmail,
			wantErr:   false,
		},
		{
			name:      "空文字（異常系）",
			plaintext: "",
			wantErr:   true,
			errMsg:    "暗号化する文字列が空です",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encrypted, err := service.Encrypt(tc.plaintext)
			if tc.wantErr {
				if err == nil {
					t.Fatal("エラーが発生することを期待しましたが、nilでした")
				}
				if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("エラーメッセージが一致しません: got %q, want to contain %q", err.Error(), tc.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("暗号化に失敗しました: %v", err)
			}

			if encrypted == "" {
				t.Fatal("暗号化結果が空です")
			}

			if encrypted == tc.plaintext {
				t.Error("暗号化結果が平文と同じです")
			}
		})
	}
}

func TestCryptoService_Decrypt(t *testing.T) {
	service, err := NewCryptoService("82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d")
	if err != nil {
		t.Fatalf("暗号化サービスの作成に失敗しました: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "有効なメールアドレス（正常系）",
			plaintext: testEmail,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encrypted, err := service.Encrypt(tc.plaintext)
			if err != nil {
				t.Fatalf("暗号化に失敗しました: %v", err)
			}

			decrypted, err := service.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("復号化に失敗しました: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("復号化結果が一致しません: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestCryptoService_Encrypt_Randomness(t *testing.T) {
	service, err := NewCryptoService("82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d")
	if err != nil {
		t.Fatalf("暗号化サービスの作成に失敗しました: %v", err)
	}

	plaintext := testEmail

	// 同じ平文を複数回暗号化して、異なる暗号文が生成されることを確認
	encrypted1, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("1回目の暗号化に失敗しました: %v", err)
	}

	encrypted2, err := service.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("2回目の暗号化に失敗しました: %v", err)
	}

	if encrypted1 == encrypted2 {
		t.Error("同じ平文を暗号化して同じ暗号文が生成されました（ランダムノンスが使われていません）")
	}

	// 両方とも正しく復号化できることを確認
	decrypted1, err := service.Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("1回目の復号化に失敗しました: %v", err)
	}

	decrypted2, err := service.Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("2回目の復号化に失敗しました: %v", err)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Errorf("復号化結果が一致しません: got %q and %q, want %q", decrypted1, decrypted2, plaintext)
	}
}

func TestCryptoService_Decrypt_InvalidInput(t *testing.T) {
	service, err := NewCryptoService("82a4d4469e7ab4916e72b05c0ee926d822ea02ed8524a73131414bd48157be9d")
	if err != nil {
		t.Fatalf("暗号化サービスの作成に失敗しました: %v", err)
	}

	tests := []struct {
		name        string
		ciphertext  string
		wantErr     bool
		errContains string
	}{
		{
			name:        "無効な16進数文字列（異常系）",
			ciphertext:  "invalid_hex",
			wantErr:     true,
			errContains: "encoding/hex",
		},
		{
			name:        "短すぎる暗号文（異常系）",
			ciphertext:  "abcd",
			wantErr:     true,
			errContains: "暗号文が短すぎます",
		},
		{
			name:        "空文字（異常系）",
			ciphertext:  "",
			wantErr:     true,
			errContains: "暗号文が短すぎます",
		},
		{
			name:        "有効な16進数だが破損したデータ（異常系）",
			ciphertext:  "1234567890abcdef1234567890abcdef12345678",
			wantErr:     true,
			errContains: "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.Decrypt(tc.ciphertext)

			if !tc.wantErr {
				if err != nil {
					t.Fatalf("予期しないエラー: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("エラーが発生することを期待しましたが、nilでした")
			}

			if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
				t.Errorf("エラーメッセージに%qが含まれていません: %q", tc.errContains, err.Error())
			}
		})
	}
}
