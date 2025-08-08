package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	userDomain "github.com/minminseo/recall-setter/domain/user"
)

func TestEmailVerificationRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name         string
		verification *userDomain.EmailVerification
		want         *userDomain.EmailVerification
		wantErr      bool
	}{
		{
			name: "既存ユーザーで認証作成（正常系）",
			verification: func() *userDomain.EmailVerification {
				ev, _ := userDomain.ReconstructEmailVerification(
					uuid.New().String(),
					"550e8400-e29b-41d4-a716-446655440004", // 既存ユーザー
					"hashed_verification_code_success",
					time.Now().Add(10 * time.Minute),
				)
				return ev
			}(),
			want: nil, // 動的に期待値を設定するためnil
			wantErr: false,
		},
		{
			name: "新規ユーザーで認証作成（外部キー制約違反）",
			verification: func() *userDomain.EmailVerification {
				ev, _ := userDomain.ReconstructEmailVerification(
					uuid.New().String(),
					uuid.New().String(),
					"hashed_verification_code_new",
					time.Now().Add(10 * time.Minute),
				)
				return ev
			}(),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewEmailVerificationRepository()

			err := repo.Create(ctx, tc.verification)

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

			// 作成された認証情報を取得して検証
			createdVerification, err := repo.FindByUserID(ctx, tc.verification.UserID())
			if err != nil {
				t.Errorf("作成された認証情報の取得に失敗: %v", err)
				return
			}

			if createdVerification == nil {
				t.Error("認証情報が作成されませんでした")
				return
			}

			// 動的に生成されるフィールドを期待値に設定するため、新しい期待値を作成
			want, _ := userDomain.ReconstructEmailVerification(
				createdVerification.ID(),
				"550e8400-e29b-41d4-a716-446655440004",
				"hashed_verification_code_success",
				createdVerification.ExpiresAt(),
			)

			// 期待値との比較
			if diff := cmp.Diff(want, createdVerification, cmp.AllowUnexported(userDomain.EmailVerification{})); diff != "" {
				t.Errorf("Create() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEmailVerificationRepository_FindByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name      string
		userID    string
		want      *userDomain.EmailVerification
		wantErr   bool
		hasRecord bool
	}{
		{
			name:   "ユーザー1の認証情報を取得（正常系）",
			userID: "550e8400-e29b-41d4-a716-446655440001",
			want: func() *userDomain.EmailVerification {
				ev, _ := userDomain.ReconstructEmailVerification(
					"c50e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"hashed_verification_code_1",
					time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
				)
				return ev
			}(),
			wantErr:   false,
			hasRecord: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewEmailVerificationRepository()

			verification, err := repo.FindByUserID(ctx, tc.userID)

			if tc.wantErr {
				if err == nil {
					t.Error("エラーが発生するはずですが、発生しませんでした")
				}
				if verification != nil {
					t.Error("認証情報がnilであるべきですが、値が返されました")
				}
				return
			}

			if err != nil && tc.hasRecord {
				t.Errorf("予期しないエラー: %v", err)
			}

			if tc.hasRecord {
				if verification == nil {
					t.Error("認証情報がnilです")
					return
				}

				// 期待値との比較
				if diff := cmp.Diff(tc.want, verification, cmp.AllowUnexported(userDomain.EmailVerification{})); diff != "" {
					t.Errorf("FindByUserID() mismatch (-want +got):\n%s", diff)
				}
			} else {
				if verification != nil {
					t.Error("認証情報が存在しないはずですが、値が返されました")
				}
			}
		})
	}
}

func TestEmailVerificationRepository_DeleteByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	PrepareTestDatabase(t)
	defer CleanupTestDatabase(t)

	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "ユーザーIDで認証情報を削除（正常系）",
			userID:  "550e8400-e29b-41d4-a716-446655440001",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := GetTestContext()
			repo := NewEmailVerificationRepository()

			err := repo.DeleteByUserID(ctx, tc.userID)

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

			// 削除後に本当に削除されたかを確認
			verification, _ := repo.FindByUserID(ctx, tc.userID)

			if verification != nil {
				t.Error("認証情報が削除されていません")
			}
		})
	}
}
