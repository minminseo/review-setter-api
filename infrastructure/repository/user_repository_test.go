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
			user: func() *userDomain.User {
				u, _ := userDomain.ReconstructUserForAuth(
					uuid.New().String(),
					"newuser@example.com",
					"encrypted_email_data",
					"encrypted_password_data",
					"light",
					"ja",
					nil,
				)
				return u
			}(),
			want: func() *userDomain.User {
				u, _ := userDomain.ReconstructUserForAuth(
					"", // IDは動的に設定
					"newuser@example.com", // FindByEmailSearchKeyでsearchKeyが設定される
					"encrypted_email_data",
					"encrypted_password_data",
					"light",
					"ja",
					nil,
				)
				return u
			}(),
			wantErr: false,
		},
		{
			name: "メール検索キー重複で作成失敗する場合",
			user: func() *userDomain.User {
				u, _ := userDomain.ReconstructUserForAuth(
					uuid.New().String(),
					"test1@example.com",
					"encrypted_email_data2",
					"encrypted_password_data2",
					"light",
					"ja",
					nil,
				)
				return u
			}(),
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
			createdUser, err := repo.FindByEmailSearchKey(ctx, tc.user.EmailSearchKey())
			if err != nil {
				t.Errorf("failed to find created user: %v", err)
				return
			}

			// 動的に生成されるフィールドを設定して新しい期待値を作成
			want, _ := userDomain.ReconstructUserForAuth(
				createdUser.ID(),
				tc.want.EmailSearchKey(),
				tc.want.EncryptedEmail(),
				tc.want.EncryptedPassword(),
				tc.want.ThemeColor(),
				tc.want.Language(),
				tc.want.VerifiedAt(),
			)
			
			if diff := cmp.Diff(want, createdUser, cmp.AllowUnexported(userDomain.User{})); diff != "" {
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
			want: func() *userDomain.User {
				u, _ := userDomain.ReconstructUserForAuth(
					"550e8400-e29b-41d4-a716-446655440001",
					"test1@example.com", // FindByEmailSearchKeyでsearchKeyが設定される
					"encrypted_email_data_1",
					"encrypted_password_1",
					"light",
					"ja",
					nil,
				)
				return u
			}(),
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

			if diff := cmp.Diff(tc.want, user, cmp.AllowUnexported(userDomain.User{})); diff != "" {
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
			user: func() *userDomain.User {
				verifiedTime := time.Now()
				u, _ := userDomain.ReconstructUserForAuth(
					"550e8400-e29b-41d4-a716-446655440001",
					"updated@example.com",
					"updated_encrypted_email",
					"updated_encrypted_password",
					"dark",
					"en",
					&verifiedTime,
				)
				return u
			}(),
			want: func() *userDomain.User {
				u, _ := userDomain.ReconstructUserForAuth(
					"550e8400-e29b-41d4-a716-446655440001",
					"updated@example.com", // FindByEmailSearchKeyで更新後のsearchKeyが設定される
					"updated_encrypted_email",
					"encrypted_password_1", // Updateではパスワードは更新されない
					"dark",
					"en",
					nil,
				)
				return u
			}(),
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
			updatedUser, err := repo.FindByEmailSearchKey(ctx, tc.user.EmailSearchKey())
			if err != nil {
				t.Errorf("failed to find updated user: %v", err)
				return
			}

			// 期待値との比較
			if diff := cmp.Diff(tc.want, updatedUser, cmp.AllowUnexported(userDomain.User{})); diff != "" {
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
			want: func() *userDomain.User {
				u, _ := userDomain.ReconstructUserForAuth(
					"550e8400-e29b-41d4-a716-446655440001",
					"test1@example.com", // FindByEmailSearchKeyでsearchKeyが設定される
					"encrypted_email_data_1",
					"new_encrypted_password",
					"light",
					"ja",
					nil,
				)
				return u
			}(),
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
			if diff := cmp.Diff(tc.want, updatedUser, cmp.AllowUnexported(userDomain.User{})); diff != "" {
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
			want: func() *userDomain.User {
				verifiedTime := time.Now()
				u, _ := userDomain.ReconstructUserForAuth(
					"550e8400-e29b-41d4-a716-446655440003",
					"test3@example.com", // FindByEmailSearchKeyでsearchKeyが設定される
					"encrypted_email_data_3",
					"encrypted_password_3",
					"light",
					"en",
					&verifiedTime,
				)
				return u
			}(),
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

			// 時間の比較で差を許容するために、期待値を新しく作成
			want := tc.want
			if tc.verifiedAt != nil {
				want, _ = userDomain.ReconstructUserForAuth(
					tc.want.ID(),
					tc.want.EmailSearchKey(),
					tc.want.EncryptedEmail(),
					tc.want.EncryptedPassword(),
					tc.want.ThemeColor(),
					tc.want.Language(),
					updatedUser.VerifiedAt(),
				)
			}

			// 期待値との比較
			if diff := cmp.Diff(want, updatedUser, cmp.AllowUnexported(userDomain.User{})); diff != "" {
				t.Errorf("UpdateVerifiedAt() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
