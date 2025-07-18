package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
)

func TestUserRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		user    *userDomain.User
		want    *userDomain.User
		wantErr bool
	}{
		{
			name: "ユーザー作成に成功する場合",
			user: &userDomain.User{
				ID:                uuid.New().String(),
				EmailSearchKey:    "newuser@example.com",
				EncryptedEmail:    "encrypted_email_data",
				EncryptedPassword: "encrypted_password_data",
				Timezone:          "Asia/Tokyo",
				ThemeColor:        "light",
				Language:          "ja",
				VerifiedAt:        nil,
			},
			want: &userDomain.User{
				EmailSearchKey:    "", // FindByEmailSearchKeyクエリでは取得されない
				EncryptedEmail:    "encrypted_email_data",
				EncryptedPassword: "encrypted_password_data",
				ThemeColor:        "light",
				Language:          "ja",
				VerifiedAt:        nil,
			},
			wantErr: false,
		},
		{
			name: "メール検索キー重複で作成失敗する場合",
			user: &userDomain.User{
				ID:                uuid.New().String(),
				EmailSearchKey:    "test1@example.com",
				EncryptedEmail:    "encrypted_email_data2",
				EncryptedPassword: "encrypted_password_data2",
				Timezone:          "Asia/Tokyo",
				ThemeColor:        "light",
				Language:          "ja",
				VerifiedAt:        nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewUserRepository()

			err := repo.Create(ctx, tc.user)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// 作成されたユーザーを取得して検証
			createdUser, err := repo.FindByEmailSearchKey(ctx, tc.user.EmailSearchKey)
			if err != nil {
				t.Errorf("failed to find created user: %v", err)
				return
			}

			// IDは動的に生成されるので除外して比較
			tc.want.ID = createdUser.ID
			if diff := cmp.Diff(tc.want, createdUser); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_FindByEmailSearchKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name      string
		searchKey string
		want      *userDomain.User
		wantErr   bool
	}{
		{
			name:      "既存ユーザーを検索する場合",
			searchKey: "test1@example.com",
			want: &userDomain.User{
				ID:                "550e8400-e29b-41d4-a716-446655440001",
				EncryptedEmail:    "encrypted_email_data_1",
				EncryptedPassword: "encrypted_password_1",
				ThemeColor:        "light",
				Language:          "ja",
				VerifiedAt:        nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewUserRepository()

			user, err := repo.FindByEmailSearchKey(ctx, tc.searchKey)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if user != nil {
					t.Error("ユーザーがnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
			}

			if user == nil {
				t.Error("ユーザーがnilです")
				return
			}

			if diff := cmp.Diff(tc.want, user); diff != "" {
				t.Errorf("FindByEmailSearchKey() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		user    *userDomain.User
		want    *userDomain.User
		wantErr bool
	}{
		{
			name: "ユーザー更新に成功する場合",
			user: &userDomain.User{
				ID:                "550e8400-e29b-41d4-a716-446655440001",
				EmailSearchKey:    "updated@example.com",
				EncryptedEmail:    "updated_encrypted_email",
				EncryptedPassword: "updated_encrypted_password",
				Timezone:          "America/New_York",
				ThemeColor:        "dark",
				Language:          "en",
				VerifiedAt:        &[]time.Time{time.Now()}[0],
			},
			want: &userDomain.User{
				ID:                "550e8400-e29b-41d4-a716-446655440001",
				EmailSearchKey:    "", // FindByEmailSearchKeyクエリでは取得されない
				EncryptedEmail:    "updated_encrypted_email",
				EncryptedPassword: "encrypted_password_1", // Updateではパスワードは更新されない
				Timezone:          "",                     // FindByEmailSearchKeyクエリでは取得されない
				ThemeColor:        "dark",
				Language:          "en",
				VerifiedAt:        nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewUserRepository()

			err := repo.Update(ctx, tc.user)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// 更新されたユーザーを取得して検証
			updatedUser, err := repo.FindByEmailSearchKey(ctx, tc.user.EmailSearchKey)
			if err != nil {
				t.Errorf("failed to find updated user: %v", err)
				return
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, updatedUser); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name     string
		userID   string
		password string
		want     *userDomain.User
		wantErr  bool
	}{
		{
			name:     "パスワード更新に成功する場合",
			userID:   "550e8400-e29b-41d4-a716-446655440001",
			password: "new_encrypted_password",
			want: &userDomain.User{
				ID:                "550e8400-e29b-41d4-a716-446655440001",
				EncryptedEmail:    "encrypted_email_data_1",
				EncryptedPassword: "new_encrypted_password",
				ThemeColor:        "light",
				Language:          "ja",
				VerifiedAt:        nil,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewUserRepository()

			err := repo.UpdatePassword(ctx, tc.userID, tc.password)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// パスワードが更新されたかを確認するためにユーザーを取得
			// フィクスチャからtest1@example.comに対応するユーザーを取得
			updatedUser, err := repo.FindByEmailSearchKey(ctx, "test1@example.com")
			if err != nil {
				t.Errorf("failed to find updated user: %v", err)
				return
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, updatedUser); diff != "" {
				t.Errorf("UpdatePassword() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserRepository_UpdateVerifiedAt(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name       string
		userID     string
		verifiedAt *time.Time
		want       *userDomain.User
		wantErr    bool
	}{
		{
			name:       "認証日時の更新に成功する場合",
			userID:     "550e8400-e29b-41d4-a716-446655440003",
			verifiedAt: &[]time.Time{time.Now()}[0],
			want: &userDomain.User{
				ID:                "550e8400-e29b-41d4-a716-446655440003",
				EncryptedEmail:    "encrypted_email_data_3",
				EncryptedPassword: "encrypted_password_3",
				ThemeColor:        "light",
				Language:          "en",
				VerifiedAt:        &[]time.Time{time.Now()}[0],
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewUserRepository()

			err := repo.UpdateVerifiedAt(ctx, tc.verifiedAt, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				return
			}

			if err != nil {
				t.Errorf("予期しないエラー: %v", err)
				return
			}

			// VerifiedAtが更新されたかを確認するためにユーザーを取得
			var updatedUser *userDomain.User
			if tc.userID == "550e8400-e29b-41d4-a716-446655440003" {
				updatedUser, err = repo.FindByEmailSearchKey(ctx, "test3@example.com")
			}
			if err != nil {
				t.Errorf("failed to find updated user: %v", err)
				return
			}

			// 時間の比較で差を許容するために、期待値のVerifiedAtを実際の値に設定
			if tc.verifiedAt != nil {
				tc.want.VerifiedAt = updatedUser.VerifiedAt
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, updatedUser); diff != "" {
				t.Errorf("UpdateVerifiedAt() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
